package eval

import (
	"bytes"
	"fmt"
	"strconv"

	"Flow/src/ast"
	"Flow/src/object"
	"Flow/src/token"
	"Flow/src/utility/slice"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	if expr, ok := node.(ast.Expression); ok {
		var err object.Object
		node, err = safeSubstituteReferences(expr, env)
		if err != nil {
			return err
		}
	}

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right, node.Token)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		return evalLetExpression(node, env)
	case *ast.IdentifierLiteral:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteralExpression:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		return applyFunction(fn, node.Arguments, env)
	case *ast.StringLiteral:
		return evalStringLiteral(node, env)
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
	}

	return nil
}

func safeSubstituteReferences(node ast.Expression, env *object.Environment) (substituted ast.Expression, err object.Object) {
	defer func() {
		if r := recover(); r != nil {
			err = object.NewEvalErrorObject("Eval: failed substituting references for %T %q, %s", node, node.String(), r)
		}
	}()

	substituted = env.SubstituteReferences(node, nil)

	return
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

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
		result = Eval(statement, env)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []ast.Expression, env *object.Environment) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.NativeFunc:
		evaluatedArgs := slice.Map(args, func(expr ast.Expression) object.Object {
			return Eval(expr, env)
		})
		return fn.Fn(evaluatedArgs...)
	default:
		return object.NewEvalErrorObject("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []ast.Expression) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, &args[paramIdx])
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

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	if node.Operator == "=" {
		return evalAssignmentExpression(node, env)
	}

	operator := node.Operator
	left := Eval(node.Left, env)
	right := Eval(node.Right, env)

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

func evalLetExpression(node *ast.LetStatement, env *object.Environment) object.Object {
	if sliceLiteral, ok := node.Value.(*ast.SliceLiteral); ok {
		c := shallowCopySliceLiteral(sliceLiteral, env)
		node.Value = c
	}
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	env.Set(node.Name.Value, &node.Value)

	return object.NULL
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	switch left := node.Left.(type) {
	case *ast.IdentifierLiteral:
		return evalAssignIdentifier(left, &node.Right, env)
	case *ast.IndexExpression:
		return evalAssignIndexExpr(left, left.Index, &node.Right, env)
	default:
		return object.NewEvalErrorObject("can't assign to give type %T", node.Left)
	}
}

func evalAssignIdentifier(identifier *ast.IdentifierLiteral, right *ast.Expression, env *object.Environment) object.Object {
	_, ok := env.Get(identifier.Value)
	if !ok {
		return object.NewEvalErrorObject(fmt.Sprintf("identifier not found: %q", identifier.Value))
	}

	env.Set(identifier.Value, right)

	return object.NULL
}

func evalAssignIndexExpr(indexExpr *ast.IndexExpression, index ast.Expression, value *ast.Expression, env *object.Environment) object.Object {
	var (
		array    *ast.ArrayLiteral
		indexInt *ast.IntegerLiteral
	)

	indexInt, ok := index.(*ast.IntegerLiteral)
	if !ok {
		return object.NewEvalErrorObject("expected integer for indexing array, got=%T", index)
	}

	if array, ok = indexExpr.Left.(*ast.ArrayLiteral); ok { // if index is used on array directly
		array.Elements[indexInt.Value] = *value

		return object.NULL
	} else if identifier, ok := indexExpr.Left.(*ast.IdentifierLiteral); ok { // if index is used on identifier referencing array
		currentValue, ok := env.Get(identifier.Value)
		if !ok {
			return object.NewEvalErrorObject(fmt.Sprintf("identifier not found: %q", identifier.Value))
		}
		array, ok = (*currentValue).(*ast.ArrayLiteral)
		if !ok {
			return object.NewEvalErrorObject(fmt.Sprintf("identifier not array, got=%T", currentValue))
		}

		array.Elements[indexInt.Value] = *value
	} else {
		return object.NewEvalErrorObject("expected array for index expression, got =%T", indexExpr.Left)
	}

	return object.NULL
}

func evalIdentifier(node *ast.IdentifierLiteral, env *object.Environment) object.Object {
	if expr, ok := env.Get(node.Value); ok {
		return Eval(*expr, env)
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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return object.NULL
	}
}

func evalStringLiteral(s *ast.StringLiteral, env *object.Environment) object.Object {
	o := object.String{}
	parts := &s.StringParts
	var out bytes.Buffer

	for {
		stringPart := parts.Value.CharacterString
		exprPart := parts.Value.Expr

		if exprPart != nil {
			val := Eval(exprPart.Expression, env)
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
	left := Eval(node.Left, env)
	idx := Eval(node.Index, env)

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

	evaluated := Eval(node.Left, env)
	array, ok := evaluated.(*object.Array)
	if !ok {
		return object.NewEvalErrorObject("indexing for type %T not implemented", node.Left)
	}

	if node.Lower != nil {
		val := Eval(*node.Lower, env)
		if intObj, ok := val.(*object.Integer); ok {
			lower = intObj
		} else {
			return object.NewEvalErrorObject("lower bound of slice must be of type integer got=%s", val.Type())
		}
	}
	if node.Upper != nil {
		val := Eval(*node.Upper, env)
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

// shallowCopySliceLiteral makes a shallow slice copy, when slicing of identifier which points to array it makes a shallow copy to slice on
func shallowCopySliceLiteral(slice *ast.SliceLiteral, env *object.Environment) *ast.SliceLiteral {
	if _, ok := slice.Left.(*ast.ArrayLiteral); ok {
		return slice
	}
	id, ok := slice.Left.(*ast.IdentifierLiteral)
	if !ok {
		panic(fmt.Sprintf("expected slice.left to be of type *ast.IdentifierLiteral got=%T", slice.Left))
	}
	identifierValue, ok := env.Get(id.Value)
	if !ok {
		panic(fmt.Sprintf("identifier %q not found, unable to slice on", id.Value))
	}
	array, ok := (*identifierValue).(*ast.ArrayLiteral)
	if !ok {
		panic(fmt.Sprintf("slice not supported for type %T", identifierValue))
	}

	newArray := ast.ArrayLiteral{
		Token:     array.Token,
		Elements:  copyArray(array.Elements),
		Generator: array.Generator,
	}
	return &ast.SliceLiteral{
		Token: slice.Token,
		Left:  &newArray,
		Lower: slice.Lower,
		Upper: slice.Upper,
	}
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
