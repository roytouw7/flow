package parser

import (
	"Flow/src/ast"
	"Flow/src/lexer"
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
	data, err := os.ReadFile("let_statements.flow")
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

	if len(program.Statements) != 3 {
		test.Fail("program.Statements dos not contain 3 statements; got =%d", program.Statements)
		return
	}

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
	data, err := os.ReadFile("return_statements.flow")
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(data))
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(test.T(), p)

	if len(program.Statements) != 3 {
		test.Fail("program.Statements does not contain 3 statements; got=%d", len(program.Statements))
	}

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

	if len(program.Statements) != 1 {
		test.T().Fatalf("program has not enough statements; got=%d", len(program.Statements))
	}
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
