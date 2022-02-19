package ast

import "Flow/src/token"

type IdentifierLiteral struct {
	Token token.Token
	Value string
}

func (i *IdentifierLiteral) expressionNode()      {}
func (i *IdentifierLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IdentifierLiteral) String() string       { return i.Value }
