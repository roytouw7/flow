package eval

import (
	"fmt"

	"Flow/src/ast"
	"Flow/src/object"
	"Flow/src/token"
)

// todo we moeten dus ook achterhalen of een identifier bestaat uit een andere identifier...
// todo iedere identifier die bestaat uit een andere identifier moet zichzelf registreren bij die andere identifier	#2.
// todo iedere identifier moet een observable zijn	#1.
// todo iedere assignment op een identifier moet de NotifyAll aanroepen met de change	#3.

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// we now recursively evaluate every node of the program
// but identifiers should be evaluated only when they are used? in for example if statements, and identifier calls?

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

	operator, left, right := node.Operator, *unwrapObservable(Eval(node.Left, env)), *unwrapObservable(Eval(node.Right, env))

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

func unwrapObservable(o object.Object) *object.Object {
	if observable, ok := o.(*object.Observable); ok {
		return unwrapObservable(observable.Value)
	}

	return &o
}

func evalLetExpression(node *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env) // todo but value can also be composition of observables? Not as easy as now
	if isError(val) {
		return val
	}

	observable := object.NewObservable(val)

	// if value is an observable we register for future changes to update our own value
	if o, ok := val.(*object.Observable); ok {
		o.Register(observable)
	}

	env.Set(node.Name.Value, observable)

	return NULL
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)

	identifier, ok := node.Left.(*ast.IdentifierLiteral)
	if !ok {
		return newEvalErrorObject("can't assign to non-identifier type, got=%T", node.Left)
	}

	val, ok := env.Get(identifier.Value)
	if !ok {
		return newEvalErrorObject(fmt.Sprintf("identifier not found: %q", identifier.Value))
	}

	if observable, ok := val.(*object.Observable); ok {
		observable.Value = right
		observable.NotifyAll(right)
	} else {
		env.Set(identifier.Value, right)
	}

	return NULL
}

func evalIdentifier(node *ast.IdentifierLiteral, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newEvalErrorObject(fmt.Sprintf("identifier not found: %s", node.Value))
	}

	return val
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
