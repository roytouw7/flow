package parser

import (
	"regexp"

	"Flow/src/ast"
	"Flow/src/helpers"
	"Flow/src/token"
)

// todo these regexps should be built from atomic pieces
const (
	arrowFnRegexp = `\(.*\)\s*=>` // whitespaces now optional because they are removed upon lexing
	groupedExprRegexp = `\(.+\)`
)

type (
	prefixParseStatementFn func() ast.Statement	// todo are these needed? Can't we use the types from the parser instead?
	infixParseStatementFn  func(left ast.Expression) ast.Statement
)

type template struct {
	match string // regex to match
	fn    any    // fn to use on match
	limit *int   // limit of characters to peek before resulting in false for template
}

// parseFnTemplateMatch matches a parse function on basis of a string match on source code
func (p *parser) parseFnTemplateMatch(prefixMatchers []template) any {
	for _, matcher := range prefixMatchers {
		tokens := make([]*token.Token, 1)
		tokens[0] = p.curToken

		if matcher.limit != nil {
			for i := 1; i < *matcher.limit; i++ {
				ok, t := p.peekTokenN(i)
				if !ok {
					break
				}

				tokens = append(tokens, t)
				if isMatch := p.tryMatchTokens(tokens, matcher.match); isMatch {
					return matcher.fn
				}
			}
		} else {
			var peekIndex = 1
			for {
				ok, t := p.peekTokenN(peekIndex)
				if !ok {
					break
				}

				tokens = append(tokens, t)
				if isMatch := p.tryMatchTokens(tokens, matcher.match); isMatch {
					return matcher.fn
				}

				peekIndex++
			}
		}
	}

	return nil
}

// tryMatchTokens peeks next token, adds to tokens slice, matches against matchString
func (p *parser) tryMatchTokens(tokens []*token.Token, matchString string) bool {
	stringRepresentation := tokensToString(tokens)

	ok, err := regexp.MatchString(matchString, stringRepresentation)
	if err != nil {
		panic(err) // todo create appropriate error
	}

	return ok
}

func tokensToString(tokens []*token.Token) string {
	return helpers.Reduce(tokens,
		func(result string, tok *token.Token) string {
			if tok.Literal != "\n" && tok.Literal != " " {
				return result + tok.Literal
			}
			return result
		}, "")
}
