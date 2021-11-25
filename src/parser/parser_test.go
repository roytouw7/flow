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
	checkParseErrors(t, p)
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
