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

// todo add more complex assignment values

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

// todo add more complex expression returns

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
