package ast

import (
	"fmt"

	"Flow/src/token"
)

type SliceLiteral struct {
	Token token.Token
	Left  *Expression
	Right *Expression
}

func (s *SliceLiteral) expressionNode() {}
func (s *SliceLiteral) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SliceLiteral) String() string {
	if s.Left != nil && s.Right != nil {
		return fmt.Sprintf("%s:%s", s.Left, s.Right)
	} else if s.Left != nil {
		return fmt.Sprintf("%s:", s.Left)
	} else if s.Right != nil {
		return fmt.Sprintf(":%s", s.Right)
	}

	return ":"
}
