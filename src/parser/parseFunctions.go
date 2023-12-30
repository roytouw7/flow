package parser

import (
	"strconv"

	"Flow/src/ast"
	"Flow/src/error"
	"Flow/src/token"
	"Flow/src/utility/convert"
	"Flow/src/utility/linkedList"
)

// todo functions like these should return the error if something goes wrong to wrap the calling context, parse function should return error
func (p *parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: *p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.registerError(cerr.Wrap(cerr.ParseIntegerLiteralError(p.curToken), "parseIntegerLiteral"))
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

func (p *parser) parseLParenExpression() ast.Expression {
	var (
		arrowFunction prefixParseFn = p.parseFunctionLiteralExpression
		groupedExpr   prefixParseFn = p.parseGroupedExpression
	)

	templates := []template{
		{
			match: arrowFnRegexp,
			fn:    arrowFunction,
			limit: nil,
		},
		{
			match: groupedExprRegexp,
			fn:    groupedExpr,
			limit: nil,
		},
	}

	parseFn := p.parseFnTemplateMatch(templates)

	if fn, ok := parseFn.(prefixParseFn); ok {
		return fn()
	}

	return nil
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

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type != token.LBRACE {
		err := cerr.UnexpectedCharError(p.peekToken, "{")
		p.registerError(cerr.Wrap(err, "parseIfExpression", "following if statement condition"))
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

func (p *parser) parseFunctionLiteralExpression() ast.Expression {
	lit := &ast.FunctionLiteralExpression{
		Token: *p.curToken,
	}

	if p.curToken.Type != token.LPAREN {
		err := cerr.UnexpectedCharError(p.curToken, token.LPAREN)
		p.registerError(cerr.Wrap(err, "parseFunctionLiteralExpression", "following function literal declaration"))
		return nil
	}

	parameters, err := p.parseFunctionParameters()
	if err != nil {
		p.registerError(cerr.Wrap(err, "parseFunctionLiteralExpression"))
		return nil
	}

	lit.Parameters = parameters

	p.nextToken()

	if p.peekToken.Type != token.ARROW {
		err := cerr.UnexpectedCharError(p.peekToken, token.ARROW)
		p.registerError(cerr.Wrap(err, "parseFunctionLiteralExpression", "following function literal parameter list declaration"))
		return nil
	}

	p.nextToken()

	if p.peekToken.Type != token.LBRACE {
		err := cerr.UnexpectedCharError(p.peekToken, token.LBRACE)
		p.registerError(cerr.Wrap(err, "parseFunctionLiteralExpression", "following function literal parameter list declaration"))
		return nil
	}

	p.nextToken()

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *parser) parseFunctionParameters() ([]*ast.IdentifierLiteral, cerr.ParseError) {
	var identifiers []*ast.IdentifierLiteral

	if p.peekToken.Type == token.RPAREN {
		return identifiers, nil
	}

	p.nextToken()

	ident := &ast.IdentifierLiteral{
		Token: *p.curToken,
		Value: p.curToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		ident := &ast.IdentifierLiteral{
			Token: *p.curToken,
			Value: p.curToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if p.peekToken.Type != token.RPAREN {
		err := cerr.UnexpectedCharError(p.peekToken, ")")
		return nil, cerr.Wrap(err, "parseFunctionParameters")
	}

	return identifiers, nil
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

func (p *parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:     *p.curToken,
		Function:  function,
		Arguments: p.parseExpressionList(token.RPAREN),
	}

	return exp
}

func (p *parser) parseStringLiteral() ast.Expression {
	exp := ast.StringLiteral{Token: *p.curToken}

	p.nextToken()

	l := linkedList.LinkedList[ast.StringLiteralPart]{}
	for p.curToken.Literal != token.STRING_DELIMITER {
		part := ast.StringLiteralPart{}
		if p.curToken.Type == token.STRING_TEMPLATE_OPEN {
			p.nextToken()
			part.Expr = p.parseExpressionStatement()
			p.nextToken()
			p.nextToken()
		}

		if p.curToken.Type == token.STRING_CHARACTERS {
			part.CharacterString = &p.curToken.Literal
			p.nextToken()
		}

		if p.curToken.Type == token.EOF {
			panic("Reached end of file before closing string!")
		}

		l.Push(part)
	}

	if p.curToken.Type != token.STRING_DELIMITER {
		panic("Not closing string!")
	}

	if l.Value == nil { // empty string literal will be set to empty string value
		l.Value = &ast.StringLiteralPart{CharacterString: convert.NewString("")}
	}

	exp.StringParts = l
	return &exp
}

func (p *parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: *p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *parser) parseLBracketExpression(left ast.Expression) ast.Expression {
	var (
		sliceParseFn infixParseFn = p.parseSliceLiteralExpression
		indexParseFn infixParseFn = p.parseIndexExpression
	)

	templates := []template{
		{
			match: sliceRegexp,
			fn:    sliceParseFn,
			limit: nil,
		},
		{
			match: indexRegexp,
			fn:    indexParseFn,
			limit: nil,
		},
	}

	parseFn := p.parseFnTemplateMatch(templates)

	if fn, ok := parseFn.(infixParseFn); ok {
		return fn(left)
	}

	return nil
}

func (p *parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: *p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RBRACKET {
		panic("Not closing array index")
	}

	p.nextToken()

	return exp
}

func (p *parser) parseSliceLiteralExpression(left ast.Expression) ast.Expression {
	exp := &ast.SliceLiteral{Token: *p.curToken, Left: left}

	if p.peekToken.Type == token.RBRACKET {
		return exp
	}

	p.nextToken()


	if p.curToken.Type != token.COLON {
		lower := p.parseExpression(SLICE)
		exp.Lower = &lower
		p.nextToken()
	}

	if p.curToken.Type != token.COLON {
		panic("No slice literal")
	}

	p.nextToken()

	if p.curToken.Type == token.RBRACKET {
		return exp
	}

	upper := p.parseExpression(SLICE)
	exp.Upper = &upper

	p.nextToken()

	if p.curToken.Type != token.RBRACKET {
		panic("Not closing array index")
	}

	return exp
}
