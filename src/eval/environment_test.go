package eval

import (
	"fmt"

	"Flow/src/ast"
	"Flow/src/object"

	"github.com/stretchr/testify/assert"
)

//func (test *Suite) TestReplaceWithOuterScopeValue() {
//	var n ast.Expression = &(ast.IntegerLiteral{Value: 4})
//	n1 := createN1Infix("n", "-", *createInteger(1))
//	n2 := createN1Infix("n", "-", *createInteger(1))
//	n3 := createN1Infix("n", "-", *createInteger(1))
//	n4 := createN1Infix("n", "-", *createInteger(1))
//	// n4 should evaluate to ((((4 - 1) - 1) - 1) - 1) => 0
//
//	env := object.NewEnvironment()
//	env.Set("n", &n)
//	env = object.NewEnclosedEnvironment(env)
//	env.Set("n", &n1)
//	env = object.NewEnclosedEnvironment(env)
//	env.Set("n", &n2)
//	env = object.NewEnclosedEnvironment(env)
//	env.Set("n", &n3)
//	env = object.NewEnclosedEnvironment(env)
//
//	substituted := env.ReplaceWithOuterScopeValue2(n4)
//	result := Eval(substituted, env)
//	testIntegerObject(test.T(), result, 0)
//}

func (test *Suite) TestSubstituteReferencesWithAssignment() {
	// x = x + y + 7 with y in an outer environment
	env := object.NewEnvironment()
	env.Set("y", createInteger(7))
	env = object.NewEnclosedEnvironment(env)
	env.Set("x", createInteger(9))

	infix := createN1Infix("y", "+", *createInteger(7)) // y + 7
	infix = createN1Infix("x", "+", infix)              // x + (y + 7)
	infix = createN1Infix("x", "=", infix)              // x = (x + (y + 7))

	result := env.SubstituteReferences(infix, nil)

	// assert result to be substituted

	infixExp, ok := result.(*ast.InfixExpression)
	if !ok {
		test.T().Errorf("expected result to be infix expression, got=%T", result)
	}
	leftIdent, ok := infixExp.Left.(*ast.IdentifierLiteral)
	if !ok {
		test.T().Errorf("expected left hand side of assignment to be identifier, got=%T", infixExp.Left)
	}
	assert.Equal(test.T(), "x", leftIdent.Value)
	assert.Equal(test.T(), "=", infixExp.Operator)

	infixExp2, ok := infixExp.Right.(*ast.InfixExpression)
	if !ok {
		assert.Fail(test.T(), "expected right hand side to be infix expression, got=%T", infixExp.Right)
	}
	int1, ok := infixExp2.Left.(*ast.IntegerLiteral)
	if !ok {
		assert.Fail(test.T(), "expected left hand side to be integer literal, got=%T", infixExp2.Left)
	}
	assert.Equal(test.T(), int64(9), int1.Value)

	infixExp3, ok := infixExp2.Right.(*ast.InfixExpression)
	if !ok {
		assert.Fail(test.T(), "expected right hand side to be infix expression, got=%T", infixExp2.Left)
	}
	identY, ok := infixExp3.Left.(*ast.IdentifierLiteral)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected left hand side to be identifier literal, got=%T", infixExp3.Left))
	}
	int3, ok := infixExp3.Right.(*ast.IntegerLiteral)
	if !ok {
		assert.Fail(test.T(), "expected right hand side to be integer literal, got=%T", infixExp3.Right)
	}

	assert.Equal(test.T(), "+", infixExp3.Operator)
	assert.Equal(test.T(), "y", identY.Value)
	assert.Equal(test.T(), int64(7), int3.Value)

	// assert original infix is not mutated

	originalInfixExp, ok := infix.(*ast.InfixExpression)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected infix to be of type infix expression, got=%T", infix))
	}
	originalIdentifier, ok := originalInfixExp.Left.(*ast.IdentifierLiteral)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected left hand side to be identifier literal, got=%T", originalInfixExp.Left))
	}
	assert.Equal(test.T(), "x", originalIdentifier.Value)

	originalInfixExp2, ok := originalInfixExp.Right.(*ast.InfixExpression)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected right hand side to be infix expression, got=%T", originalInfixExp.Right))
	}
	originalIdentifier1, ok := originalInfixExp2.Left.(*ast.IdentifierLiteral)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected left hand side to be identifier literal, got=%T", originalInfixExp2.Left))
	}
	assert.Equal(test.T(), "x", originalIdentifier1.Value)
}

func createInteger(n int64) *ast.Expression {
	var obj ast.Expression = &ast.IntegerLiteral{Value: n}
	return &obj
}

// createN1Infix returns infix in form of n-1
func createN1Infix(identifier, operator string, expression ast.Expression) ast.Expression {
	var n1 ast.Expression = &(ast.InfixExpression{
		Left:     &(ast.IdentifierLiteral{Value: identifier}),
		Operator: operator,
		Right:    expression,
	})

	return n1
}
