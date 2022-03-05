package parser2

import (
	"Flow/src/ast"
	"Flow/src/token"
	"fmt"
)

// Parser input collection of *token.Token
// Parser output *ast.Program consisting of a []ast.Statement
// fetch next token
// if EOF return program
// parseStatement
// if known statement parse it
// else try parsing it as expression
// getPrefixParseFunction
// parsePrefixExpression
// leftExp = parsed prefix expression
// peek next token
// if next has infixParseFunction
// parse next as infixParseFunction
// leftExp = parsedInfix; repeat until precedence of peeked token is higher; everytime making leftExp the result of the parsedInfix expression
// else return leftExp

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Lexer interface {
	NextToken() *token.Token
}

type Parser interface {
	ParseProgram() *ast.Program
	Errors() []error
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)

type parser struct {
	l Lexer

	errors    []error
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

	// Set current and peek token
	p.nextToken()
	p.nextToken()

	return p
}

func (p *parser) Errors() []error {
	return nil
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

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *parser) registerError(err error) {
	p.errors = append(p.errors, err)
}

func unexpectedToken(expected token.Type, token *token.Token) error {
	return fmt.Errorf("expected next token to be %s, got %s as %d:%d", expected, token.Type, token.Line, token.Pos)
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

	if p.peekToken.Type == token.IDENT {
		p.nextToken()
	} else {
		p.registerError(unexpectedToken(token.IDENT, p.peekToken))
		return nil
	}

	stmt.Name = &ast.IdentifierLiteral{Token: *p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
	} else {
		p.registerError(unexpectedToken(token.ASSIGN, p.peekToken))
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *parser) parseReturnStatement() *ast.ReturnStatement {
	return nil
}

func (p *parser) parseExpressionStatement() *ast.ExpressionStatement {
	return nil
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.registerError(fmt.Errorf("no parse function found for token type %s", p.curToken.Type))
		return nil
	}

	leftExp := prefix()

	// todo check if it should be parse as infix expression and parse it if so

	return leftExp

}
