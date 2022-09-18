package parser2

import (
	"encoding/json"
	"os"
	"testing"

	"Flow/src/ast"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestLetStatements() {
	program := createProgramFromFile(test.T(), "test_assets/let_statements.flow", 5)

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

func (test *Suite) TestReturnStatements() {
	program := createProgramFromFile(test.T(), "test_assets/return_statements.flow", 5)

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
	program := createProgramFromFile(test.T(), "test_assets/identifier_expressions.flow", 3)

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
	program := createProgramFromFile(test.T(), "test_assets/integer_literal_expressions.flow", 3)

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
	program := createProgramFromFile(test.T(), "test_assets/boolean_literal_expressions.flow", 2)

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
	program := createProgramFromFile(test.T(), "test_assets/prefix_expressions.flow", 4)

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

func (test *Suite) TestInfixExpressions() {
	program := createProgramFromFile(test.T(), "test_assets/infix_expressions.flow", 11)

	tests := []struct {
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{leftValue: 5, operator: "+", rightValue: 5},
		{leftValue: 5, operator: "-", rightValue: 5},
		{leftValue: 5, operator: "*", rightValue: 5},
		{leftValue: 5, operator: "/", rightValue: 5},
		{leftValue: 5, operator: ">", rightValue: 5},
		{leftValue: 5, operator: "<", rightValue: 5},
		{leftValue: 5, operator: "==", rightValue: 5},
		{leftValue: 5, operator: "!=", rightValue: 5},
		{leftValue: true, operator: "==", rightValue: true},
		{leftValue: true, operator: "!=", rightValue: false},
		{leftValue: false, operator: "==", rightValue: false},
	}

	for i, tt := range tests {
		stmt, ok := program.Statements[i].(*ast.ExpressionStatement)

		if !ok {
			test.T().Errorf("program.Sttements[%d] is not ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testInfixExpression(test.T(), stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

// todo complete test; if this one works it should be on the same level as the book again but refactored
func (test *Suite) TestOperatorPrecedenceParsing() {
	data, err := os.ReadFile("test_assets/precedence_tests.json")
	if err != nil {
		panic(err)
	}

	var tests []struct {
		Input    string
		Expected string
	}

	if err = json.Unmarshal(data, &tests); err != nil {
		panic(err)
	}

	for _, tt := range tests {
		createProgram(test.T(), tt.Input, 1)
	}
}
