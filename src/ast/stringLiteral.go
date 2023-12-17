package ast

import (
	"bytes"

	"Flow/src/token"
	"Flow/src/utility/linkedList"
)

type StringLiteralPart struct {
	CharacterString *string
	Expr            *ExpressionStatement
}

type StringLiteral struct {
	Token       token.Token
	StringParts linkedList.LinkedList[StringLiteralPart]
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string {
	var out bytes.Buffer

	link := s.StringParts

	for {
		stringLiteralPart := link.Value
		str, expr := stringLiteralPart.CharacterString, stringLiteralPart.Expr

		if str != nil {
			out.WriteString(*str)
		}

		if expr != nil {
			out.WriteString(expr.String())
		}

		if link.HasNext() {
			link = *link.Next()
		} else {
			break
		}
	}

	return out.String()
}
