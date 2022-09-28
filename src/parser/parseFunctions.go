package parser

import (
	"fmt"
	"strconv"

	"Flow/src/ast"
	"Flow/src/token"
)

func (p *parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: *p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.registerError(fmt.Errorf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *parser) parseIdentifier() ast.Expression {
	return &ast.IdentifierLiteral{Token: *p.curToken, Value: p.curToken.Literal}
}

func (p *parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: *p.curToken,
		Value: p.curToken.Type == token.TRUE,
	}
}

// todo add error logging
func (p *parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return exp
	}

	return nil
}

func (p *parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    *p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    *p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := precedences[p.curToken.Type]
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{
		Token: *p.curToken,
	}

	if p.peekToken.Type != token.LPAREN {
		p.registerError(fmt.Errorf("expected ( after if statement"))
		return nil
	}

	p.nextTokenN(2)
	expression.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		p.registerError(fmt.Errorf("expected ( after if statement"))
		return nil
	}

	p.nextToken()

	if p.peekToken.Type != token.LBRACE {
		p.registerError(fmt.Errorf("expected { to open if statement body after condition declaration"))
		return nil
	}

	p.nextToken()

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if p.peekToken.Type != token.LBRACE {
			return nil
		}

		p.nextToken()
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *parser) parseTernaryExpression(left ast.Expression) ast.Expression {
	expression := &ast.TernaryExpression{
		Token:     *p.curToken,
		Condition: left,
	}

	p.nextToken()

	expression.Consequence = p.parseExpression(TERNARY)

	p.nextTokenN(2)

	expression.Alternative = p.parseExpression(TERNARY)

	return expression
}

func (p *parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: *p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}
