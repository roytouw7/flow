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

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
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
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
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
		return newEvalErrorObject("not a function: %s", fn.Type())
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
		return newEvalErrorObject("%sunknown operator: %s%s", tokenToPos(token), operator, right.Type())
	}
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	if node.Operator == "=" {
		return evalAssignmentExpression(node, env)
	}

	operator, left, right := node.Operator, *unwrapObservable(Eval(node.Left, env), env), *unwrapObservable(Eval(node.Right, env), env)

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right) // Only works with booleans because bool objects are reused so memory address matches
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newEvalErrorObject("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newEvalErrorObject("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func unwrapObservable(o object.Object, env *object.Environment) *object.Object {
	if observable, ok := o.(*object.Observable); ok {
		val := Eval(*observable.Value, env)
		return unwrapObservable(val, env)
	}

	return &o
}

func evalLetExpression(node *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	observable := object.NewObservable(&node.Value)

	// if value is an observable we register for future changes to update our own value
	if o, ok := val.(*object.Observable); ok {
		o.Register(observable)
	}

	env.Set(node.Name.Value, &node.Value)

	return NULL
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	identifier, ok := node.Left.(*ast.IdentifierLiteral)
	if !ok {
		return newEvalErrorObject("can't assign to non-identifier type, got=%T", node.Left)
	}

	expr, ok := env.Get(identifier.Value)
	if !ok {
		return newEvalErrorObject(fmt.Sprintf("identifier not found: %q", identifier.Value))
	}

	val := Eval(*expr, env)

	if observable, ok := val.(*object.Observable); ok {
		observable.Value = &node.Right
		observable.NotifyAll(&node.Right)
	} else {
		substituted := substituteSelfReference(node.Right, identifier.Value, expr)
		env.Set(identifier.Value, &substituted)
	}

	return NULL
}

// substituteSelfReference evaluation results in stack overflow if assignment is self referencing like a = a;
func substituteSelfReference(node ast.Expression, self string, value *ast.Expression) ast.Expression {
	switch node := node.(type) {
	case *ast.IdentifierLiteral:
		if node.Value == self { // identifier is self referencing, substitute for the value
			return *value
		}
		return node
	case *ast.PrefixExpression:
		right := substituteSelfReference(node.Right, self, value)
		node.Right = right
		return node
	case *ast.InfixExpression:
		left := substituteSelfReference(node.Left, self, value)
		right := substituteSelfReference(node.Right, self, value)
		node.Left = left
		node.Right = right
		return node
	default:
		return node
	}
}

func evalIdentifier(node *ast.IdentifierLiteral, env *object.Environment) object.Object {
	if expr, ok := env.Get(node.Value); ok {
		return Eval(*expr, env)
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newEvalErrorObject(fmt.Sprintf("identifier not found: %s", node.Value))
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, tok token.Token) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newEvalErrorObject("%sunknown operator: -%s", tokenToPos(tok), right.Type())
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
		return newEvalErrorObject("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
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
			val := Eval(exprPart, env)
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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newEvalErrorObject("index operator not support for array literal: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
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

func newEvalErrorObject(format string, a ...interface{}) *object.EvalError {
	return &object.EvalError{Message: fmt.Sprintf(format, a...)}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
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
