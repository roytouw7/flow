package parser

import (
	"Flow/src/ast"
	"Flow/src/lexer"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestLetStatements() {
	data, err := os.ReadFile("test_assets/let_statements.flow")
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(data))
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(test.T(), p)
	if program == nil {
		test.Fail("ParseProgram() returned nil")
		return
	}

	test.checkProgramLines(program, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !test.testLetStatement(stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func (test *Suite) testLetStatement(s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		test.T().Errorf("s.TokenLiteral not 'let'; got=%T", s)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		test.T().Errorf("s not *ast.LetStatement; got=%T", s)
	}

	if letStmt.Name.Value != name {
		test.T().Errorf("letStmt.Name.Value not '%s'; got =%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		test.T().Errorf("s.Name not '%s'; got=%s", name, letStmt.Name)
	}

	return true
}

func (test *Suite) TestReturnStatements() {
	data, err := os.ReadFile("test_assets/return_statements.flow")
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(data))
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(test.T(), p)
	test.checkProgramLines(program, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			test.T().Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			test.T().Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func (test *Suite) TestIdentifierExpression() {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(test.T(), p)

	test.checkProgramLines(program, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.T().Fatalf("program.Statements[0] is not ast.ExpresionStatement; got=%T", program.Statements[0])
	}

	testLiteralExpression(test.T(), stmt.Expression, "foobar")
}

func (test *Suite) TestIntegerLiteralExpression() {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(test.T(), p)

	test.checkProgramLines(program, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.T().Fatalf("program.Statements[0] is not ast.ExpresionStatement; got=%T", program.Statements[0])
	}

	testLiteralExpression(test.T(), stmt.Expression, 5)
}

func (test *Suite) TestBooleanLiteralExpression() {
	input := "false;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(test.T(), p)

	test.checkProgramLines(program, 1)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		test.T().Fatalf("program.Statements[0] is not ast.ExpresionStatement; got=%T", program.Statements[0])
	}

	testLiteralExpression(test.T(), stmt.Expression, false)
}

func (test *Suite) TestParsingPrefixExpressions() {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(test.T(), p)

		test.checkProgramLines(program, 1)
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			test.T().Fatalf("program.Statements[0] is not ast.ExpresionStatement; got=%T", program.Statements[0])
		}
		if testPrefixExpression(test.T(), stmt.Expression, tt.operator, tt.value) {
			return
		}
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	default:
		t.Errorf("type of expression not handled, got=%T", exp)
		return false
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral; got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d; got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d; got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	literal, ok := bl.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.BooleanLiteral; got=%T", bl)
		return false
	}
	if literal.Value != value {
		t.Errorf("literal.Value not %t, got=%T", true, literal.Value)
		return false
	}
	if literal.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("literal.TokenLiteral not %t; got=%s", value, literal.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'; got=%q", operator, opExp.Operator)
		return false
	}

	return true
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.PrefixExpression)

	if !ok {
		t.Errorf("exp is not ast.PrefixExpression, got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'; got=%q", operator, opExp.Operator)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, il ast.Expression, value string) bool {
	ident, ok := il.(*ast.IdentifierLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral; got=%T", il)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s; got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("integ.TokenLiteral not %s; got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func (test *Suite) TestParsingInfixExpressions() {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{input: "5 + 5", leftValue: 5, operator: "+", rightValue: 5},
		{input: "5 - 5", leftValue: 5, operator: "-", rightValue: 5},
		{input: "5 * 5", leftValue: 5, operator: "*", rightValue: 5},
		{input: "5 / 5", leftValue: 5, operator: "/", rightValue: 5},
		{input: "5 > 5", leftValue: 5, operator: ">", rightValue: 5},
		{input: "5 < 5", leftValue: 5, operator: "<", rightValue: 5},
		{input: "5 == 5", leftValue: 5, operator: "==", rightValue: 5},
		{input: "5 != 5", leftValue: 5, operator: "!=", rightValue: 5},
		{input: "true == true", leftValue: true, operator: "==", rightValue: true},
		{input: "true != false", leftValue: true, operator: "!=", rightValue: false},
		{input: "false == false", leftValue: false, operator: "==", rightValue: false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(test.T(), p)
		test.checkProgramLines(program, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			test.T().Fatalf("program.Statements[0] is not ast.ExpressionStatement; got=%T", program.Statements[0])
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
		l := lexer.New(tt.Input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(test.T(), p)

		actual := program.String()
		if actual != tt.Expected {
			test.T().Fatalf("expected=%q; got=%q", tt.Expected, actual)
		}
	}
}

func checkParseErrors(t *testing.T, p Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err)
	}
	t.FailNow()
}

func (test *Suite) checkProgramLines(p *ast.Program, expectedLines int) {
	if len(p.Statements) != expectedLines {
		test.T().Fatalf("program does not have the correct amount of   statements; got=%d expected=%d", len(p.Statements), expectedLines)
	}
}
