package ast

import (
	"bytes"
	"strings"

	"Flow/src/token"
	"Flow/src/utility/slice"
)

type ArrayLiteral struct {
	Token     token.Token
	Elements  []Expression
	Generator bool
	// or something like x for x in 1...9 if expr else expr then expr
}

func (a *ArrayLiteral) expressionNode() {}

func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := slice.Map(a.Elements, func(el Expression) string {
		return el.String()
	})

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
