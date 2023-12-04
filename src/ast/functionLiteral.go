package ast

import (
	"bytes"
	"strings"

	"Flow/src/token"
)

type FunctionLiteralExpression struct {
	Token      token.Token
	Parameters []*IdentifierLiteral
	Body       *BlockStatement
}

func (fl *FunctionLiteralExpression) expressionNode()      {}
func (fl *FunctionLiteralExpression) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteralExpression) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}
