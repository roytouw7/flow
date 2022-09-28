package ast

import (
	"bytes"

	"Flow/src/token"
)

type TernaryExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }

func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString(te.Condition.String())
	out.WriteString("?")
	out.WriteString(te.Consequence.String())
	out.WriteString(":")
	out.WriteString(te.Alternative.String())

	return out.String()
}
