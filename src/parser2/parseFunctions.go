package parser2

import (
	"Flow/src/ast"
	"Flow/src/token"
	"fmt"
	"strconv"
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
