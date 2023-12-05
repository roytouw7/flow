package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"Flow/src/ast"
	"Flow/src/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestLetStatements() {
	program := CreateProgramFromFile(test.T(), "test_assets/let_statements.flow", 5)

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
	program := CreateProgramFromFile(test.T(), "test_assets/return_statements.flow", 5)

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
	program := CreateProgramFromFile(test.T(), "test_assets/identifier_expressions.flow", 3)

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
	program := CreateProgramFromFile(test.T(), "test_assets/integer_literal_expressions.flow", 3)

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
	program := CreateProgramFromFile(test.T(), "test_assets/boolean_literal_expressions.flow", 2)

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
	program := CreateProgramFromFile(test.T(), "test_assets/prefix_expressions.flow", 4)

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
	program := CreateProgramFromFile(test.T(), "test_assets/infix_expressions.flow", 11)

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
		CreateProgram(test.T(), tt.Input, 1)
	}
}

func (test *Suite) TestIfExpressions() {
	program := CreateProgramFromFile(test.T(), "test_assets/if_expressions.flow", 2)

	var tests = []struct {
		condition   string
		consequence []string
	}{
		{
			condition:   "(x < y)",
			consequence: []string{"x"},
		}, {
			condition:   "(x > y)",
			consequence: []string{"x", "y"},
		},
	}

	for i, tt := range tests {
		es, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			test.T().Errorf("statement %d is no *ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testIfExpression(test.T(), es.Expression, tt.condition, tt.consequence)
	}
}

func (test *Suite) TestIfElseExpressions() {
	program := CreateProgramFromFile(test.T(), "test_assets/if_else_expressions.flow", 2)

	var tests = []struct {
		condition   string
		consequence []string
		alternative []string
	}{
		{
			condition:   "(x < y)",
			consequence: []string{"alfa"},
			alternative: []string{"beta"},
		},
		{
			condition:   "(x > y)",
			consequence: []string{"alfa", "beta"},
			alternative: []string{"gamma"},
		},
	}

	var counter int
	for i, tt := range tests {
		es, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			test.T().Errorf("statement %d is no *ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testIfElseExpression(test.T(), es.Expression, tt.condition, tt.consequence, tt.alternative)
		counter = i + 1
	}

	if counter != len(program.Statements) {
		test.T().Errorf("not all program statements tested, expected=%d got=%d", len(program.Statements), counter)
	}
}

func (test *Suite) TestTernaryExpressions() {
	program := CreateProgramFromFile(test.T(), "test_assets/ternary_expressions.flow", 2)

	var tests = []struct {
		condition   string
		consequence string
		alternative string
	}{
		{
			condition:   "(a > b)",
			consequence: "(a + 1)",
			alternative: "(b + 2)",
		}, {
			condition:   "true",
			consequence: "false",
			alternative: "true",
		},
	}

	var counter int
	for i, tt := range tests {
		es, ok := program.Statements[i].(*ast.ExpressionStatement)
		if !ok {
			test.T().Errorf("statement %d is no *ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testTernaryExpression(test.T(), es.Expression, tt.condition, tt.consequence, tt.alternative)
		counter = i + 1
	}

	if counter != len(program.Statements) {
		test.T().Errorf("not all program statements tested, expected=%d got=%d", len(program.Statements), counter)
	}
}

func (test *Suite) TestGroupedExpressions() {
	CreateProgramFromFile(test.T(), "test_assets/grouped_expressions.flow", 1)
}

func (test *Suite) TestFunctionLiteralExpressions() {
	program := CreateProgramFromFile(test.T(), "test_assets/function_literal_expressions.flow", 2)

	var tests = []struct {
		parameters []string
		statements []string
	}{
		{
			parameters: []string{},
			statements: []string{"return 7;"},
		},
		{
			parameters: []string{"a", "b"},
			statements: []string{"return (a * b);"},
		},
	}

	// todo a lot of repeating logic in these tests yet, and all the custom errors...
	var counter int
	for i, tt := range tests {
		es, ok := program.Statements[i].(*ast.LetStatement)
		if !ok {
			test.T().Errorf("statement %d is no *ast.ExpressionStatement; got=%T", i, program.Statements[i])
		}

		testFunctionLiteralExpression(test.T(), es.Value, tt.parameters, tt.statements)
		counter = i + 1
	}

	if counter != len(program.Statements) {
		test.T().Errorf("not all program statements tested, expected=%d got=%d", len(program.Statements), counter)
	}
}

func (test *Suite) TestParseTemplateMatcher() {
	var parseFn prefixParseStatementFn

	p := createParserFromFile("test_assets/sample_program.flow")

	parseFn = func() ast.Statement {
		stmt := &ast.ExpressionStatement{
			Token:      token.Token{},
			Expression: nil,
		}
		return stmt
	}

	var limit = 10

	input := []template{
		{
			match: `\(.*\)\s*=>`, // whitespace is stripped by parser now, but anything except for whitespace between is invalid
			fn:    parseFn,
			limit: &limit,
		},
	}

	p.nextTokenN(3)

	result := p.parseFnTemplateMatch(input)

	prefixParseFn, ok := result.(prefixParseStatementFn)
	if !ok {
		assert.Failf(test.T(), "testParseTemplateMatcher test", "function is type %t, expected prefixParseFn", result)
	}

	assert.NotNil(test.T(), prefixParseFn)
	assert.Equal(test.T(), fmt.Sprintf("%p", parseFn), fmt.Sprintf("%p", prefixParseFn))
}

func (test *Suite) TestParseTemplateMatcher_TwoTemplates() {
	var (
		parseFnArrowFn     prefixParseStatementFn
		parseFnGroupedExpr infixParseStatementFn
	)

	p := createParserFromFile("test_assets/grouped_expressions.flow")

	parseFnArrowFn = func() ast.Statement {
		stmt := &ast.ExpressionStatement{
			Token:      token.Token{},
			Expression: nil,
		}
		return stmt
	}

	parseFnGroupedExpr = func(left ast.Expression) ast.Statement {
		stmt := &ast.ExpressionStatement{
			Token:      token.Token{},
			Expression: nil,
		}
		return stmt
	}

	var limit = 10

	input := []template{
		{
			match: `\(.*\)\s*=>`,
			fn:    parseFnArrowFn,
			limit: nil,
		},
		{
			match: `\(.+\)`,
			fn:    parseFnGroupedExpr,
			limit: &limit,
		},
	}

	p.nextTokenN(3)

	result := p.parseFnTemplateMatch(input)

	prefixParseFn, ok := result.(infixParseStatementFn)
	if !ok {
		assert.Failf(test.T(), "testParseTemplateMatcher test", "function is type %t, expected prefixParseFn", result)
	}

	assert.NotNil(test.T(), prefixParseFn)
	assert.Equal(test.T(), fmt.Sprintf("%p", parseFnGroupedExpr), fmt.Sprintf("%p", prefixParseFn))
}

func (test *Suite) TestParseTemplateMatcher_NoMatch() {
	var parseFn prefixParseStatementFn

	p := createParserFromFile("test_assets/sample_program.flow")

	parseFn = func() ast.Statement {
		stmt := &ast.ExpressionStatement{
			Token:      token.Token{},
			Expression: nil,
		}
		return stmt
	}

	var limit = 10

	input := []template{
		{
			match: `\[\]\[\]`, // not present in sample program
			fn:    parseFn,
			limit: &limit,
		},
	}

	p.nextTokenN(3)

	result := p.parseFnTemplateMatch(input)

	assert.Nil(test.T(), result)
}

func (test *Suite) TestCallExpressionParsing() {
	program := CreateProgramFromFile(test.T(), "test_assets/call_expressions.flow", 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.Failf("TestCallExpressionParsing", "stmt is not ast.ExpressionsStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		test.Failf("TestCallExpressionParsing", "expr is not ast.CallExpression, got=%T", stmt.Expression)
	}

	if !testIdentifier(test.T(), exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		test.Failf("TestCallExpressionParsing", "number of arguments does not equal 3, got=%d", len(exp.Arguments))
	}
}
