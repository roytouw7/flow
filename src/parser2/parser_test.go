package parser2

import (
	"Flow/src/ast"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// todo add checkTestExecution function to check if all tests completed
func (test *Suite) TestLetStatements() {
	program := createProgram(test.T(), "test_assets/let_statements.flow", 5)

	tests := []struct {
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
		{"foo", "bar"},
		{"flag", true},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(test.T(), stmt, tt.expectedIdentifier, tt.expectedValue)
	}
}

// todo add checkTestExecution function to check if all tests completed
func (test *Suite) TestReturnStatements() {
	program := createProgram(test.T(), "test_assets/return_statements.flow", 5)

	tests := []struct {
		expectedReturnValue interface{}
	}{
		{5},
		{10},
		{993322},
		{"foobar"},
		{false},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testReturnStatement(test.T(), stmt, tt.expectedReturnValue)
	}
}

// todo add checkTestExecution function to check if all tests completed
func (test *Suite) TestIdentifierExpression() {
	program := createProgram(test.T(), "test_assets/identifier_expressions.flow", 3)

	tests := []struct {
		expectedIdentifier interface{}
	}{
		{"foobar"},
		{"django"},
		{"lara777"},
	}

	for i, tt := range tests {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)

		if !ok {
			test.T().Errorf("program.Sttements[%d] is not ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testLiteralExpression(test.T(), stmt.Expression, tt.expectedIdentifier)
	}
}

// todo add checkTestExecution function to check if all tests completed
func (test *Suite) TestIntegerLiteralExpression() {
	program := createProgram(test.T(), "test_assets/integer_literal_expressions.flow", 3)

	tests := []struct {
		expectedReturnValue interface{}
	}{
		{1337},
		{10},
		{993322},
	}

	for i, tt := range tests {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)

		if !ok {
			test.T().Errorf("program.Sttements[%d] is not ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testLiteralExpression(test.T(), stmt.Expression, tt.expectedReturnValue)
	}
}

// todo add checkTestExecution function to check if all tests completed
func (test *Suite) TestBooleanLiteralExpression() {
	program := createProgram(test.T(), "test_assets/boolean_literal_expressions.flow", 2)

	tests := []struct {
		expectedReturnValue interface{}
	}{
		{false},
		{true},
	}

	for i, tt := range tests {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)

		if !ok {
			test.T().Errorf("program.Sttements[%d] is not ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testLiteralExpression(test.T(), stmt.Expression, tt.expectedReturnValue)
	}
}

func (test *Suite) TestPrefixExpressions() {
	program := createProgram(test.T(), "test_assets/prefix_expressions.flow", 4)

	tests := []struct {
		expectedOperator string
		expectedValue    interface{}
	}{
		{"!", 5},
		{"-", 9},
		{"!", true},
		{"-", "foo"},
	}

	for i, tt := range tests {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)

		if !ok {
			test.T().Errorf("program.Sttements[%d] is not ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testPrefixExpression(test.T(), stmt.Expression, tt.expectedOperator, tt.expectedValue)
	}
}
