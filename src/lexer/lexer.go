package lexer

import (
	"Flow/src/iterator"
	"Flow/src/token"
	"unicode"
)

type lexer struct {
	iterator iterator.StringIterator
}

func New(input string) *lexer {
	i := iterator.New(input)
	l := &lexer{iterator: i}
	return l
}

func (l *lexer) NextToken() *token.Token {
	if !l.iterator.HasNext() {
		return createEOFSymbolToken()
	}

	ch, meta, err := l.getNextNonWhiteSpaceCharacter()
	if err != nil {
		panic(err)
	}

	ok, tok := l.isSymbolToken(ch, meta)
	if ok {
		return tok
	}

	ok, tok = l.isStringLiteral(ch, meta)
	if ok {
		identifierType := token.LookupIdentType(tok.Literal)
		tok.Type = identifierType
		return tok
	}

	ok, tok = l.isNumericLiteralToken(ch, meta)
	if ok {
		return tok
	}

	return createIllegalSymbolToken(meta)
}

func createIllegalSymbolToken(meta *iterator.MetaData) *token.Token {
	return token.New(token.ILLEGAL, "???", meta.RelPos, meta.Line)
}

func createEOFSymbolToken() *token.Token {
	return token.New(token.EOF, token.EOF, -1, -1)
}

func (l *lexer) getNextNonWhiteSpaceCharacter() (rune, *iterator.MetaData, error) {
	for {
		ch, meta, err := l.iterator.Next()
		if err != nil {
			return ch, meta, err
		}

		if !(ch != '\n' && unicode.IsSpace(ch)) {
			return ch, meta, err
		}
	}
}

// todo default rune type is 0, EOL should be given different symbol in iterator to differentiate

func (l *lexer) isSymbolToken(ch rune, meta *iterator.MetaData) (bool, *token.Token) {
	newToken := curriedSymbolTokenConstructor(meta)

	switch ch {
	case '=':
		switch {
		case l.isMultiSymbolToken('='):
			return true, newToken(token.EQ)
		default:
			return true, newToken(token.ASSIGN)
		}
	case '+':
		return true, newToken(token.PLUS)
	case '-':
		return true, newToken(token.MINUS)
	case '!':
		switch {
		case l.isMultiSymbolToken('='):
			return true, newToken(token.NOT_EQ)
		default:
			return true, newToken(token.BANG)
		}
	case '*':
		return true, newToken(token.ASTERISK)
	case '/':
		return true, newToken(token.SLASH)
	case '<':
		return true, newToken(token.LT)
	case '>':
		return true, newToken(token.GT)
	case ',':
		return true, newToken(token.COMMA)
	case ';':
		return true, newToken(token.SEMICOLON)
	case '\n':
		return true, newToken(token.NEWLINE)
	case '(':
		return true, newToken(token.LPAREN)
	case ')':
		return true, newToken(token.RPAREN)
	case '{':
		return true, newToken(token.LBRACE)
	case '}':
		return true, newToken(token.RBRACE)
	default:
		return false, newToken(token.UNKNOWN)
	}
}

func (l *lexer) isStringLiteral(ch rune, meta *iterator.MetaData) (bool, *token.Token) {
	if !unicode.IsLetter(ch) {
		return false, token.NewSymbol(token.UNKNOWN, meta.RelPos, meta.Line)
	}

	literal := make([]rune, 0)
	literal = append(literal, ch)

	l.appendLiteralUntil(&literal, func(ch rune) bool {
		return !(unicode.IsLetter(ch) || unicode.IsDigit(ch))
	})

	return true, token.New(token.IDENT, string(literal), meta.RelPos, meta.Line)
}

func (l *lexer) isNumericLiteralToken(ch rune, meta *iterator.MetaData) (bool, *token.Token) {
	if !unicode.IsDigit(ch) {
		return false, token.NewSymbol(token.UNKNOWN, meta.RelPos, meta.Line)
	}

	literal := make([]rune, 0)
	literal = append(literal, ch)

	l.appendLiteralUntil(&literal, func(ch rune) bool {
		return !unicode.IsDigit(ch)
	})

	return true, token.New(token.INT, string(literal), meta.RelPos, meta.Line)
}

func (l *lexer) appendLiteralUntil(literal *[]rune, delimitFn func(ch rune) bool) {
	for l.iterator.HasNext() {
		p, err := l.iterator.Peek()
		if err != nil {
			panic(err)
		}

		if delimitFn(p) {
			break
		}

		next, _, err := l.iterator.Next()
		if err != nil {
			panic(err)
		}

		*literal = append(*literal, next)
	}
}

func curriedSymbolTokenConstructor(meta *iterator.MetaData) func(t token.Type) *token.Token {
	return func(t token.Type) *token.Token {
		return token.NewSymbol(t, meta.RelPos, meta.Line)
	}
}

func (l *lexer) isMultiSymbolToken(chs ...rune) bool {
	for i, ch := range chs {
		p, err := l.iterator.PeekN(i + 1)
		if err != nil {
			panic(err)
		}
		if p != ch {
			return false
		}

		if _, _, err = l.iterator.Next(); err != nil {
			panic(err)
		}
	}

	return true
}
