package ast

import (
	"bytes"

	"Flow/src/token"
)

type TernaryExpression struct {
	Token     token.Token
	Condition Expression
	TrueExp   *BlockStatement
	FalseExp  *BlockStatement
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }

func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString(te.Condition.String())
	out.WriteString("?")
	out.WriteString(te.TrueExp.String())
	out.WriteString(":")
	out.WriteString(te.FalseExp.String())

	return out.String()
}
