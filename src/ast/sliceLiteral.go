package ast

import (
	"fmt"

	"Flow/src/token"
)

type SliceLiteral struct {
	Token token.Token
	Left  Expression
	Lower *Expression
	Upper *Expression
}

func (s *SliceLiteral) expressionNode() {}
func (s *SliceLiteral) TokenLiteral() string {
	return s.Token.Literal
}
func (s *SliceLiteral) String() string {
	var (
		lowerStr string
		upperStr string
	)

	if s.Lower != nil {
		lowerStr = (*(*s).Lower).String()

	}
	if s.Upper != nil {
		upperStr = (*(*s).Upper).String()
	}

	if s.Lower != nil && s.Upper != nil {
		return fmt.Sprintf("(%s[%s:%s])", s.Left.String(), lowerStr, upperStr)
	} else if s.Lower != nil {
		return fmt.Sprintf("(%s[%s:])", s.Left.String(), lowerStr)
	} else if s.Upper != nil {
		return fmt.Sprintf("(%s[:%s])", s.Left.String(), upperStr)
	}

	return fmt.Sprintf("(%s[:])", s.Left.String())
}
