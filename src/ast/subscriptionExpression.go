package ast

import (
	"bytes"

	"Flow/src/token"
)

// source can be a pipeline or a direct value
// pipeline can have a source of pipeline or direct value too

type SubscriptionExpression struct {
	Token  token.Token
	Source Expression	// always an identifier?
	Body   Expression	// either function literal or identifier
}

func (s *SubscriptionExpression) expressionNode() {}

func (s *SubscriptionExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *SubscriptionExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(s.Source.String())
	out.WriteString(") ~> { ")

	out.WriteString(s.Body.String())

	out.WriteString("}")
	return out.String()
}
