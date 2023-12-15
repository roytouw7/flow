package lexer

import (
	"os"
	"testing"

	cerr "Flow/src/error"
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
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.ARROW, "=>"},
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
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "199"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.STRING_DELIMITER, "\""},
		{token.STRING_CHARACTERS, "foobar"},
		{token.STRING_DELIMITER, "\""},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			err := cerr.TestUnexpectedTokenError("TestPeekN", i, tok, tt.expectedType)
			test.T().Error(err)
		}

		if tok.Literal != tt.expectedLiteral {
			err := cerr.TestUnexpectedValueFor("TestPeekN", i, "literal", tok.Literal, tt.expectedLiteral)
			test.T().Error(err)
		}
	}
}

func (test *Suite) TestPeekN() {
	data, err := os.ReadFile("test_assets/test_program.flow")
	test.Nil(err)

	var input = string(data)
	l := New(input)

	type expected struct {
		ok    bool
		token *token.Token
	}

	var tests = []struct {
		input  int
		output expected
	}{
		{
			input:  0,
			output: expected{false, nil},
		},
		{
			input:  1,
			output: expected{true, token.New(token.LET, "let", 1, 1)},
		},
		{
			input:  2,
			output: expected{true, token.New(token.IDENT, "five", 5, 1)},
		},
		{
			input:  3,
			output: expected{true, token.New(token.ASSIGN, "=", 10, 1)},
		},
		{
			input:  4,
			output: expected{true, token.New(token.INT, "5", 12, 1)},
		},
		{
			input:  10,
			output: expected{true, token.New(token.INT, "10", 11, 2)},
		},
		{
			input:  20,
			output: expected{true, token.New(token.IDENT, "y", 15, 4)},
		},
		{
			input:  30,
			output: expected{true, token.New(token.RBRACE, "}", 1, 6)},
		},
		{
			input:  250,
			output: expected{false, nil},
		},
	}

	for i, tt := range tests {
		ok, tok := l.PeekN(tt.input)
		if tt.output.ok {
			if !assert.Equal(test.T(), tt.output.token, tok) {
				err := cerr.TestUnexpectedTokenError("TestPeekN", i, tok, tt.output.token.Type)
				test.T().Error(err)
			}
			if ok != tt.output.ok {
				err := cerr.TestUnexpectedValueFor("TestPeekN", i, "ok", ok, tt.output.ok)
				test.T().Error(err)
			}
		} else {
			if ok != false {
				err := cerr.TestUnexpectedValueFor("TestPeekN", i, "ok", ok, false)
				test.T().Error(err)
			}
		}
	}
}

func (test *Suite) TestEatString() {
	l := New("\"foobar\";")
	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.STRING_DELIMITER, "\""},
		{token.STRING_CHARACTERS, "foobar"},
		{token.STRING_DELIMITER, "\""},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	for _, tt := range tests {
		tok := l.NextToken()
		test.Equal(tt.expectedToken, tok.Type)
		test.Equal(tt.expectedLiteral, tok.Literal)
	}
}

func (test *Suite) TestEatString_WithInterpolation() {
	l := New("\"foo ${7 + 9} bar\";")
	tests := []struct {
		expectedToken   token.Type
		expectedLiteral string
	}{
		{token.STRING_DELIMITER, "\""},
		{token.STRING_CHARACTERS, "foo "},
		{token.STRING_TEMPLATE_OPEN, "${"},
		{token.INT, "7"},
		{token.PLUS, "+"},
		{token.INT, "9"},
		{token.RBRACE, "}"},
		{token.STRING_CHARACTERS, " bar"},
		{token.STRING_DELIMITER, "\""},
		{token.SEMICOLON, ";"},
		{token.EOF, "EOF"},
	}

	for _, tt := range tests {
		tok := l.NextToken()
		test.Equal(tt.expectedToken, tok.Type)
		test.Equal(tt.expectedLiteral, tok.Literal)
	}
}

func (test *Suite) TestEOF() {
	l := New("0")

	tt := []token.Type{token.INT, token.EOF, token.EOF, token.EOF}

	for _, t := range tt {
		result := l.NextToken()
		if result.Type != t {
			test.T().Errorf("expected token type %s got %T", t, result)
		}
	}
}
