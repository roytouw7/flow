package cerr

import (
	"fmt"

	"Flow/src/token"
)

type AError interface {
	compileErrorInterface
	aError() // to discriminate ParseError from other Errors when required
}

type ParseError interface {
	compileErrorInterface
	parseError() // to discriminate ParseError from other Errors when required
}

type parseError struct {
	compileError
}

func (p *parseError) parseError() {}

type ParseFnType string

const (
	Prefix ParseFnType = "prefix"
	Infix  ParseFnType = "infix"
)

func UnexpectedCharError(tok *token.Token, expected string) ParseError {
	msg := fmt.Sprintf("expected character %q, got %q instead", expected, tok.Literal)
	return &parseError{
		compileError{
			msg, tok,
		},
	}
}

func MissingParseFnError(tok *token.Token, kind ParseFnType) ParseError {
	msg := fmt.Sprintf("no %s parse function found for token %q", kind, tok.Literal)
	return &parseError{
		compileError{
			msg, tok,
		},
	}
}

func UnexpectedTokenError(tok *token.Token, expected token.Type, actual token.Type) ParseError {
	msg := fmt.Sprintf("expected token to be %q, got %q instead", expected, actual)
	return &parseError{
		compileError{
			msg, tok,
		},
	}
}

func ParseIntegerLiteralError(tok *token.Token) ParseError {
	msg := fmt.Sprintf("could not parse %q as integer", tok.Literal)
	return &parseError{
		compileError{
			msg, tok,
		},
	}
}
