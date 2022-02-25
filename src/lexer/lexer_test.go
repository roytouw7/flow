package lexer

import (
	"Flow/src/token"
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

func (test *Suite) TestSymbolToken() {
	l := New("== ") // todo remove trailing space when fixed
	t := l.NextToken()
	test.NotNil(t)

	if t.Literal != "==" {
		test.T().Errorf("expected %s, got %s", "==", t.Literal)
	}
	if t.Type != token.EQ {
		test.T().Errorf("expected %s, got %s", token.EQ, t.Type)
	}
}

func (test *Suite) TestStringLiteral() {
	l := New("al  ") // todo remove trailing space when fixed
	t := l.NextToken()
	test.NotNil(t)

	if t.Literal != "al" {
		test.T().Errorf("expected %s, got %s", "al", t.Literal)
	}
}

// TestIntEqual case caused bug before
func (test *Suite) TestIntEqual() {
	l := New("10 == 10;")

	var tests = []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		test.Equalf(tok.Type, tt.expectedType, "tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		test.Equalf(tok.Literal, tt.expectedLiteral, "tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)

	}
}

func (test *Suite) TestNextToken() {
	data, err := os.ReadFile("test_assets/test_program.flow")
	test.Nil(err)

	var input = string(data)

	var tests = []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		//{token.SEMICOLON, ";"},
		//{token.NEWLINE, "\n"},
		//{token.INT, "10"},
		//{token.NOT_EQ, "!="},
		//{token.INT, "199"},
		//{token.SEMICOLON, ";"},
		//{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		test.Equalf(tok.Type, tt.expectedType, "tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		test.Equalf(tok.Literal, tt.expectedLiteral, "tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)

	}
}
