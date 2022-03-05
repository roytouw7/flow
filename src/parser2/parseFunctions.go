package parser2

import (
	"Flow/src/ast"
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
