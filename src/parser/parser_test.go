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
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		test.T().Fatalf("exp not *ast.Identifier; got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		test.T().Errorf("ident.Value not %s; got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		test.T().Errorf("ident.TokenLiteral not %s; got=%s", "foobar", ident.TokenLiteral())
	}
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
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		test.T().Fatalf("exp not *ast.IntegerLiteral; got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		test.T().Errorf("literal.Value not %d; got=%T", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		test.T().Errorf("literal.TokenLiteral not %s; got=%s", "5", literal.TokenLiteral())
	}
}

func (test *Suite) TestParsingPrefixExpressions() {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15;", "-", 15},
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
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			test.T().Fatalf("stmt is not ast.PrefixExpressions; goy=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			test.T().Fatalf("exp.Operator is not '%s'; got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(test.T(), exp.Right, tt.integerValue) {
			return
		}
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

func (test *Suite) TestParsingInfixExpressions() {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{input: "5 + 5", leftValue: 5, operator: "+", rightValue: 5},
		{input: "5 - 5", leftValue: 5, operator: "-", rightValue: 5},
		{input: "5 * 5", leftValue: 5, operator: "*", rightValue: 5},
		{input: "5 / 5", leftValue: 5, operator: "/", rightValue: 5},
		{input: "5 > 5", leftValue: 5, operator: ">", rightValue: 5},
		{input: "5 < 5", leftValue: 5, operator: "<", rightValue: 5},
		{input: "5 == 5", leftValue: 5, operator: "==", rightValue: 5},
		{input: "5 != 5", leftValue: 5, operator: "!=", rightValue: 5},
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
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			test.T().Fatalf("stmt.Expressions is not ast.InfixExpressions; got=%T", stmt.Expression)
		}
		if !testIntegerLiteral(test.T(), exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			test.T().Fatalf("exp.Operator is not '%s'; got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(test.T(), exp.Right, tt.rightValue) {
			return
		}
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

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.errors

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
