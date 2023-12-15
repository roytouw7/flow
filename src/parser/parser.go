package parser

//todo should set up wire?

import (
	"Flow/src/ast"
	"Flow/src/error"
	"Flow/src/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGNMENT
	TERNARY
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.Type]int{
	token.QUESTION: TERNARY,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.ASSIGN:   ASSIGNMENT,
}

type Lexer interface {
	NextToken() *token.Token
	PeekN(n int) (bool, *token.Token)
}

type Parser interface {
	ParseProgram() *ast.Program
	Errors() []cerr.ParseError
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)

type parser struct {
	l Lexer // todo why isn't this embeded?

	errors    []cerr.ParseError
	curToken  *token.Token
	peekToken *token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l Lexer) Parser {
	p := &parser{
		l: l,
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.TRUE] = p.parseBooleanLiteral
	p.prefixParseFns[token.FALSE] = p.parseBooleanLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.IF] = p.parseIfExpression
	p.prefixParseFns[token.LPAREN] = p.parseLParenExpression
	p.prefixParseFns[token.STRING_DELIMITER] = p.parseStringLiteral

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.ASSIGN] = p.parseInfixExpression
	p.infixParseFns[token.QUESTION] = p.parseTernaryExpression
	p.infixParseFns[token.LPAREN] = p.parseCallExpression

	// Set current and peek token
	p.nextToken()
	p.nextToken()

	return p
}

func (p *parser) Errors() []cerr.ParseError {
	return p.errors
}

func (p *parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *parser) nextTokenN(n int) {
	for i := 0; i < n; i++ {
		p.nextToken()
	}
}

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *parser) peekTokenN(n int) (bool, *token.Token) {
	if n < 0 {
		return false, nil
	}

	if n == 1 {
		return true, p.peekToken
	}

	return p.l.PeekN(n - 1)
}

func (p *parser) logOnFailure(fn func(t token.Type) bool, t token.Type, err cerr.ParseError) bool { // todo this function only clutters the flow now, not clear what it does
	if ok := fn(t); ok {
		return true
	}

	p.registerError(err)
	return false
}

func (p *parser) incrementOnMatch(t token.Type) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}

	return false
}

func (p *parser) registerError(err cerr.ParseError) {
	p.errors = append(p.errors, err)
}

func (p *parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.NEWLINE:
		return nil
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: *p.curToken}

	if !p.logOnFailure(p.incrementOnMatch, token.IDENT, cerr.UnexpectedTokenError(p.peekToken, token.IDENT)) {
		return nil
	}

	stmt.Name = &ast.IdentifierLiteral{Token: *p.curToken, Value: p.curToken.Literal}

	if !p.logOnFailure(p.incrementOnMatch, token.ASSIGN, cerr.UnexpectedTokenError(p.peekToken, token.ASSIGN)) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	p.incrementOnMatch(token.SEMICOLON)

	return stmt
}

func (p *parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: *p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	p.incrementOnMatch(token.SEMICOLON)

	return stmt
}

func (p *parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: *p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// if ternary expression the statement token should be of the ternary expression instead of the token of the (partial)condition
	// e.g. for expression: "a > b ? a : b;" token should be set to ? instead of IDENT
	if stmt.Expression != nil && stmt.Expression.TokenLiteral() == token.QUESTION {
		tok := stmt.Token
		stmt.Token = *token.New(token.QUESTION, token.QUESTION, tok.Pos, tok.Line)
	}

	p.incrementOnMatch(token.SEMICOLON)

	return stmt
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.registerError(cerr.Wrap(cerr.MissingParseFnError(p.curToken, cerr.Prefix), "parseExpression"))
		return nil
	}

	leftExp := prefix()

	for (p.peekToken.Type != token.SEMICOLON && p.peekToken.Type != token.RBRACE) && precedence < precedences[p.peekToken.Type] {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			p.registerError(cerr.Wrap(cerr.MissingParseFnError(p.curToken, cerr.Infix), "parseExpression"))
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp

}
