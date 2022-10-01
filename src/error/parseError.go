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
	msg := fmt.Sprintf("expected character %s, got %s instead", expected, tok.Literal) // TODO maybe %q instead of %s
	return &parseError{
		compileError{
			msg, tok,
		},
	}
}

func MissingParseFnError(tok *token.Token, kind ParseFnType) ParseError {
	msg := fmt.Sprintf("no %s parse function found for token %s", kind, tok.Literal) // todo maybe use %q here to
	return &parseError{
		compileError{
			msg, tok,
		},
	}
}

func UnexpectedTokenError(tok *token.Token, expected token.Type, actual token.Type) ParseError {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", expected, actual) // todo same %q
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
