package eval

import (
	"bytes"
	"fmt"
	"strconv"

	"Flow/src/ast"
	"Flow/src/object"
	"Flow/src/token"
	"Flow/src/utility/observer"
	"Flow/src/utility/slice"
)

func Eval(node ast.Node, parent observer.Observer, env *object.Environment) object.Object {
	observable := observer.WrapNodeWithObservable(node, observer.New())
	observable, err := safeSubstituteReferences(observable, env)
	if err != nil {
		return err
	}
	if parent != nil {
		observable.Register(parent)
	}

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, nil, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, nil, env)
		return evalPrefixExpression(node.Operator, right, node.Token)
	case *object.InfixExpressionObservable:
		return evalInfixExpression(observable, env)
	case *ast.InfixExpression:
		return evalInfixExpression(observable, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IfExpression:
		return evalIfExpression(node, parent, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, nil, env)
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		return evalLetExpression(observable, env)
	case *ast.IdentifierLiteral:
		return evalIdentifier(node, parent, env)
	case *ast.FunctionLiteralExpression:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		fn := Eval(node.Function, nil, env)
		if isError(fn) {
			return fn
		}
		return applyFunction(object.WrapObjectWithObservable(fn, observer.New()), node.Arguments, env)
	case *ast.StringLiteral:
		return evalStringLiteral(node, parent, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.SliceLiteral:
		return evalSliceExpression(node, env)
	case *ast.SubscriptionExpression:
		return evalSubscriptionExpression(observable, env)
	}

	return nil
}

func safeSubstituteReferences(observable observer.ObservableNode[ast.Node], env *object.Environment) (substituted observer.ObservableNode[ast.Node], err object.Object) {
	node := observable.Node

	defer func() {
		if r := recover(); r != nil {
			err = object.NewEvalErrorObject("Eval: failed substituting references for %T %q, %s", node, node.String(), r)
		}
	}()

	substituted = env.SubstituteReferences(observable, nil)

	return
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, nil, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.EvalError:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, nil, env)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, nil, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(observable object.ObservableObject[object.Object], args []ast.Expression, env *object.Environment) object.Object {
	fn := observable.Object

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, &observable, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.NativeFunc:
		evaluatedArgs := slice.Map(args, func(expr ast.Expression) object.Object {
			return Eval(expr, &observable, env)
		})
		return fn.Fn(evaluatedArgs...)
	default:
		return object.NewEvalErrorObject("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []ast.Expression) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, observer.WrapNodeWithObservable[ast.Node](args[paramIdx], observer.New()))
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalPrefixExpression(operator string, right object.Object, token token.Token) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, token)
	default:
		return object.NewEvalErrorObject("%sunknown operator: %s%s", tokenToPos(token), operator, right.Type())
	}
}

func evalInfixExpression(observable observer.ObservableNode[ast.Node], env *object.Environment) object.Object {
	node, ok := observable.Node.(*object.InfixExpressionObservable)
	if !ok {
		panic(fmt.Sprintf("Node not infix expression, got=%T!", observable.Node))
	}

	node.Left.Register(observable)
	node.Right.Register(observable)

	observable.SetHandler(func(id observer.TraceId) {
		//fmt.Println("infix notify")
	})

	if node.Operator == "=" {
		return evalAssignmentExpression(node, observable, env)
	}

	operator := node.Operator
	left := Eval(node.Left.Node, &observable, env)
	right := Eval(node.Right.Node, &observable, env)

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right) // Only works with booleans because bool objects are reused so memory address matches
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return object.NewEvalErrorObject("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return object.NewEvalErrorObject("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalLetExpression(observable observer.ObservableNode[ast.Node], env *object.Environment) object.Object {
	node, ok := observable.Node.(*ast.LetStatement)
	if !ok {
		panic(fmt.Sprintf("Node not let statement, got=%T!", observable.Node))
	}

	//if sliceLiteral, ok := node.Value.(*ast.SliceLiteral); ok {
	//	c := shallowCopySliceLiteral(sliceLiteral, env)
	//	node.Value = c
	//}
	val := Eval(node.Value, observable, env)
	if isError(val) {
		return val
	}

	right := observer.WrapNodeWithObservable[ast.Node](node.Value, observer.New())
	right.Register(observable)

	env.Set(node.Name.Value, right)

	observable.SetHandler(func(id observer.TraceId) {
		latest, ok := env.Get(node.Name.Value)
		if !ok {
			panic(fmt.Sprintf("identifier %q is removed from closure", node.Name.Value))
		}
		//fmt.Println("letExpression notify")
		latest.Notify(&id)
	})

	return object.NULL
}

func evalAssignmentExpression(node *object.InfixExpressionObservable, parent observer.Observer, env *object.Environment) object.Object {
	right := node.Right.Node.(ast.Expression)
	switch left := node.Left.Node.(type) {
	case *ast.IdentifierLiteral:
		return evalAssignIdentifier(left, &right, env)
	//case *ast.IndexExpression:
	//	return evalAssignIndexExpr(left, left.Index, &node.Right, env)
	default:
		return object.NewEvalErrorObject("can't assign to give type %T", node.Left)
	}
}

func evalAssignIdentifier(identifier *ast.IdentifierLiteral, right *ast.Expression, env *object.Environment) object.Object {
	val, ok := env.Get(identifier.Value)
	if !ok {
		return object.NewEvalErrorObject(fmt.Sprintf("identifier not found: %q", identifier.Value))
	}

	exp := observer.WrapNodeWithObservable[ast.Node](*right, observer.New())
	exp.Register(val)

	exp.SetHandler(func(id observer.TraceId) {
		latest, ok := env.Get(identifier.Value)
		if !ok {
			panic("asd")
		}
		latest.Notify(&id)
		//fmt.Println("assign notify")
	})

	env.Set(identifier.Value, exp)
	val.Notify(nil)

	return object.NULL
}

//func evalAssignIndexExpr(indexExpr *ast.IndexExpression, index ast.Expression, value *ast.Expression, env *object.Environment) object.Object {
//	var (
//		array    *ast.ArrayLiteral
//		indexInt *ast.IntegerLiteral
//	)
//
//	indexInt, ok := index.(*ast.IntegerLiteral)
//	if !ok {
//		return object.NewEvalErrorObject("expected integer for indexing array, got=%T", index)
//	}
//
//	if array, ok = indexExpr.Left.(*ast.ArrayLiteral); ok { // if index is used on array directly
//		array.Elements[indexInt.Value] = *value
//
//		return object.NULL
//	} else if identifier, ok := indexExpr.Left.(*ast.IdentifierLiteral); ok { // if index is used on identifier referencing array
//		currentValue, ok := env.Get(identifier.Value)
//		if !ok {
//			return object.NewEvalErrorObject(fmt.Sprintf("identifier not found: %q", identifier.Value))
//		}
//		array, ok = (*currentValue).(*ast.ArrayLiteral)
//		if !ok {
//			return object.NewEvalErrorObject(fmt.Sprintf("identifier not array, got=%T", currentValue))
//		}
//
//		array.Elements[indexInt.Value] = *value
//	} else {
//		return object.NewEvalErrorObject("expected array for index expression, got =%T", indexExpr.Left)
//	}
//
//	return object.NULL
//}

func evalIdentifier(node *ast.IdentifierLiteral, parent observer.Observer, env *object.Environment) object.Object {
	if expr, ok := env.Get(node.Value); ok {
		return Eval(expr.Node, parent, env)
	}

	if builtin, ok := object.Builtins[node.Value]; ok {
		return builtin
	}

	return object.NewEvalErrorObject(fmt.Sprintf("identifier not found: %s", node.Value))
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, tok token.Token) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return object.NewEvalErrorObject("%sunknown operator: -%s", tokenToPos(tok), right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return object.NewEvalErrorObject("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, parent observer.Observer, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, parent, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, parent, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, parent, env)
	} else {
		return object.NULL
	}
}

func evalStringLiteral(s *ast.StringLiteral, parent observer.Observer, env *object.Environment) object.Object {
	o := object.String{}
	parts := &s.StringParts
	var out bytes.Buffer

	for {
		stringPart := parts.Value.CharacterString
		exprPart := parts.Value.Expr

		if exprPart != nil {
			val := Eval(exprPart.Expression, parent, env)
			out.WriteString(toString(val))
		}

		if stringPart != nil {
			out.WriteString(*stringPart)
		}

		if !parts.HasNext() {
			break
		}
		parts = parts.Next()
	}

	o.Value = out.String()
	return &o
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, nil, env)
	idx := Eval(node.Index, nil, env)

	if array, ok := left.(*object.Array); ok {
		if intIdx, ok := idx.(*object.Integer); ok {
			max := int64(len(array.Elements) - 1)
			if intIdx.Value < 0 || intIdx.Value > max {
				return object.NULL
			}
			return array.Elements[intIdx.Value]
		}
	}

	return object.NewEvalErrorObject("indexing for type %T not implemented", node.Left)
}

func evalSliceExpression(node *ast.SliceLiteral, env *object.Environment) object.Object {
	var (
		lower *object.Integer
		upper *object.Integer
	)

	evaluated := Eval(node.Left, nil, env)
	array, ok := evaluated.(*object.Array)
	if !ok {
		return object.NewEvalErrorObject("indexing for type %T not implemented", node.Left)
	}

	if node.Lower != nil {
		val := Eval(*node.Lower, nil, env)
		if intObj, ok := val.(*object.Integer); ok {
			lower = intObj
		} else {
			return object.NewEvalErrorObject("lower bound of slice must be of type integer got=%s", val.Type())
		}
	}
	if node.Upper != nil {
		val := Eval(*node.Upper, nil, env)
		if intObj, ok := val.(*object.Integer); ok {
			upper = intObj
		} else {
			return object.NewEvalErrorObject("lower bound of slice must be of type integer got=%s", val.Type())
		}
	}

	if len(array.Elements) == 0 {
		return &object.Array{}
	} else if lower != nil && upper != nil {
		return &object.Array{Elements: array.Elements[lower.Value:upper.Value]}
	} else if lower != nil {
		return &object.Array{Elements: array.Elements[lower.Value:]}
	} else if upper != nil {
		return &object.Array{Elements: array.Elements[:upper.Value]}
	}

	return &object.Array{Elements: array.Elements}

}

// copyArray makes a shallow copy of the array
func copyArray[T any](source []T) []T {
	result := make([]T, len(source))
	copy(result, source)
	return result
}

//// shallowCopySliceLiteral makes a shallow slice copy, when slicing of identifier which points to array it makes a shallow copy to slice on
//func shallowCopySliceLiteral(slice *ast.SliceLiteral, env *object.Environment) *ast.SliceLiteral {
//	if _, ok := slice.Left.(*ast.ArrayLiteral); ok {
//		return slice
//	}
//	id, ok := slice.Left.(*ast.IdentifierLiteral)
//	if !ok {
//		panic(fmt.Sprintf("expected slice.left to be of type *ast.IdentifierLiteral got=%T", slice.Left))
//	}
//	identifierValue, ok := env.Get(id.Value)
//	if !ok {
//		panic(fmt.Sprintf("identifier %q not found, unable to slice on", id.Value))
//	}
//	array, ok := (*identifierValue).(*ast.ArrayLiteral)
//	if !ok {
//		panic(fmt.Sprintf("slice not supported for type %T", identifierValue))
//	}
//
//	newArray := ast.ArrayLiteral{
//		Token:     array.Token,
//		Elements:  copyArray(array.Elements),
//		Generator: array.Generator,
//	}
//	return &ast.SliceLiteral{
//		Token: slice.Token,
//		Left:  &newArray,
//		Lower: slice.Lower,
//		Upper: slice.Upper,
//	}
//}

func evalSubscriptionExpression(observable observer.ObservableNode[ast.Node], env *object.Environment) object.Object {
	node, ok := observable.Node.(*ast.SubscriptionExpression)
	if !ok {
		panic(fmt.Sprintf("Node not SubscriptionExpression, got=%T", observable.Node))
	}

	left := node.Source
	id, ok := left.(*ast.IdentifierLiteral)
	if !ok {
		panic(fmt.Sprintf("Left hand side not IdentifierLiteral, got=%T", left))
	}

	val, ok := env.Get(id.Value)
	if !ok {
		panic(fmt.Sprintf("Identifier not found in closure %q", id.Value))
	}

	val.Register(observable)

	right := node.Body
	fnId, ok := right.(*ast.IdentifierLiteral)
	if !ok {
		panic(fmt.Sprintf("expected identifier got=%T", node.Body))
	}

	val.SetHandler(func(traceId observer.TraceId) {
		latest, ok := env.Get(id.Value)
		if !ok {
			panic(fmt.Sprintf("identifier no longer in closure %q", id.Value))
		}

		var callExpression ast.Node = &ast.CallExpression{
			Function: fnId,
			Arguments: []ast.Expression{
				latest.Node.(ast.Expression),
			},
		}
		Eval(callExpression, nil, env)
		//fmt.Println("subscription notify")
	})

	return &object.Null{}
}

func toString(obj object.Object) string {
	switch obj := obj.(type) {
	case *object.String:
		return obj.Value
	case *object.Integer:
		return strconv.FormatInt(obj.Value, 10)
	case *object.Boolean:
		if obj.Value {
			return "true"
		}
		return "false"
	case *object.Null: // todo identifier type
		return "NULL"
	}
	panic(fmt.Sprintf("can't stringify type %T", obj))
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}

	return object.FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func tokenToPos(tok token.Token) string {
	return fmt.Sprintf("%d:%d: ", tok.Line, tok.Pos)
}
