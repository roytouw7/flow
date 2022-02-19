package ast

import (
	"Flow/src/token"
	"bytes"
)

type LetStatement struct {
	Token token.Token
	Name  *IdentifierLiteral
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	// TODO is this required now semicolons are optional?
	out.WriteString(";")

	return out.String()
}
