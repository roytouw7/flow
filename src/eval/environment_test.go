package eval

import (
	"Flow/src/ast"
	"Flow/src/object"
)

func (test *Suite) TestReplaceWithOuterScopeValue() {
	var n ast.Expression = &(ast.IntegerLiteral{Value: 4})
	n1 := createN1Infix("n", "-", 1)
	n2 := createN1Infix("n", "-", 1)
	n3 := createN1Infix("n", "-", 1)
	n4 := createN1Infix("n", "-", 1)
	// n4 should evaluate to ((((4 - 1) - 1) - 1) - 1) => 0

	env := object.NewEnvironment()
	env.Set("n", &n)
	env = object.NewEnclosedEnvironment(env)
	env.Set("n", &n1)
	env = object.NewEnclosedEnvironment(env)
	env.Set("n", &n2)
	env = object.NewEnclosedEnvironment(env)
	env.Set("n", &n3)

	substituted := env.ReplaceWithOuterScopeValue2(n4)
	result := Eval(substituted, env)
	testIntegerObject(test.T(), result, 0)
}

// createN1Infix returns infix in form of n-1
func createN1Infix(identifier, operator string, delta int64) ast.Expression {
	var n1 ast.Expression = &(ast.InfixExpression{
		Left:     &(ast.IdentifierLiteral{Value: identifier}),
		Operator: operator,
		Right:    &(ast.IntegerLiteral{Value: delta}),
	})

	return n1
}
