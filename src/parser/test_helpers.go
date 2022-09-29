package parser

import (
	"fmt"
	"os"
	"testing"

	"Flow/src/ast"
	"Flow/src/lexer"
)

func createProgramFromFile(t *testing.T, fileName string, expectedStatements int) *ast.Program {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return createProgram(t, string(data), expectedStatements)
}

// TODO test cases should log which input line triggered the failing test somehow
func createProgram(t *testing.T, input string, expectedStatements int) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
		return nil
	}

	checkParseErrors(t, p)
	checkProgramLines(t, program, expectedStatements)

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
		fmt.Println(p.String())
		t.Fatalf("program does not have the correct amount of statements; got=%d expected=%d", len(p.Statements), expectedLines)
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

func testPrefixExpression(t *testing.T, s ast.Expression, operator string, value interface{}) bool {
	exp, ok := s.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("statement not *ast.PrefixExpression; got=%T", s)
		return false
	}

	if exp.Operator != operator {
		t.Errorf("expected expression operator to be %s, got %s", operator, exp.Operator)
		return false
	}

	if !testLiteralExpression(t, exp.Right, value) {
		t.Errorf("unexpected expression value")
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, leftValue interface{}, operator string, rightValue interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, leftValue) {
		return false
	}

	if !testLiteralExpression(t, opExp.Right, rightValue) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'; got=%q", operator, opExp.Operator)
		return false
	}

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
		t.Errorf("literal.Value not %t, got=%t", value, literal.Value)
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

func testIfExpression(t *testing.T, ie ast.Expression, condition string, consequenceStatements []string) bool {
	ifExpression, ok := ie.(*ast.IfExpression)
	if !ok {
		t.Errorf("ie not *ast.IfExpression; got=%T", ie)
		return false
	}

	if ifExpression.Condition.String() != condition {
		t.Errorf("expected condition %s for expression; got=%s", condition, ifExpression.Condition.String())
		return false
	}

	if ifExpression.Alternative != nil && len(ifExpression.Alternative.Statements) != 0 {
		t.Errorf("expexted no alternative blockstatement in if expression")
		return false
	}

	if !testBlockStatement(t, ifExpression.Consequence, consequenceStatements) {
		return false
	}

	return true
}

func testIfElseExpression(t *testing.T, ie ast.Expression, condition string, consequenceStatements []string, alternativeStatements []string) bool {
	ifExpression, ok := ie.(*ast.IfExpression)
	if !ok {
		t.Errorf("ie not *ast.IfExpression; got=%T", ie)
		return false
	}

	if ifExpression.Condition.String() != condition {
		t.Errorf("expected condition %s for expression; got=%s", condition, ifExpression.Condition.String())
		return false
	}

	if !testBlockStatement(t, ifExpression.Consequence, consequenceStatements) {
		return false
	}

	if !testBlockStatement(t, ifExpression.Alternative, alternativeStatements) {
		return false
	}

	return true
}

func testTernaryExpression(t *testing.T, e ast.Expression, condition string, consequence string, alternative string) bool {
	ternaryExpression, ok := e.(*ast.TernaryExpression)
	if !ok {
		t.Errorf("te not *ast.TernaryExpression; got=%T", e)
		return false
	}

	if ternaryExpression.Condition.String() != condition {
		t.Errorf("expected condition %s for expression; got=%s", condition, ternaryExpression.Condition.String())
		return false
	}

	if ternaryExpression.Consequence.String() != consequence {
		t.Errorf("expected consequence %s for expression; got=%s", consequence, ternaryExpression.Consequence.String())
		return false
	}

	if ternaryExpression.Alternative.String() != alternative {
		t.Errorf("expected alternative %s for expression; got=%s", alternative, ternaryExpression.Alternative.String())
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, bs ast.Statement, expectedStatements []string) bool {
	blockStatement, ok := bs.(*ast.BlockStatement)
	if !ok {
		t.Errorf("bs not *ast.BlockStatement; got=%T", bs)
		return false
	}

	if len(blockStatement.Statements) != len(expectedStatements) {
		t.Errorf("expected %d statements in block statements; got=%d", len(blockStatement.Statements), len(expectedStatements))
		return false
	}

	for i, stmt := range blockStatement.Statements {
		if stmt.String() != expectedStatements[i] {
			t.Errorf("expected statement %s; got=%s at index %d of BlockStatement", expectedStatements[i], stmt.String(), i)
			return false
		}
	}

	return true
}
