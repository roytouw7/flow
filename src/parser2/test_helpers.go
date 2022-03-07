package parser2

import (
	"Flow/src/ast"
	"Flow/src/lexer"
	"fmt"
	"os"
	"testing"
)

func createProgram(t *testing.T, fileName string, expectedLines int) *ast.Program {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(data))
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
		return nil
	}

	checkParseErrors(t, p)
	checkProgramLines(t, program, expectedLines)

	return program
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

func checkProgramLines(t *testing.T, p *ast.Program, expectedLines int) {
	if len(p.Statements) != expectedLines {
		t.Fatalf("program does not have the correct amount of   statements; got=%d expected=%d", len(p.Statements), expectedLines)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string, value interface{}) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'; got=%T", s)
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement; got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'; got =%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'; got=%s", name, letStmt.Name)
		return false
	}

	testLiteralExpression(t, letStmt.Value, value)

	return true
}

func testReturnStatement(t *testing.T, s ast.Statement, value interface{}) bool {
	if s.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral not 'let'; got=%T", s)
		return false
	}

	returnStmt, ok := s.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("s not *ast.ReturnStatement; got=%T", s)
		return false
	}

	testLiteralExpression(t, returnStmt.ReturnValue, value)

	return true
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

func testIdentifier(t *testing.T, il ast.Expression, value string) bool {
	ident, ok := il.(*ast.IdentifierLiteral)
	if !ok {
		t.Errorf("il not *ast.IdentifierLiteral; got=%T", il)
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
