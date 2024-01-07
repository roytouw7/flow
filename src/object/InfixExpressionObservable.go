package object

import (
	"bytes"

	"Flow/src/ast"
	"Flow/src/token"
	"Flow/src/utility/observer"
)

type InfixExpressionObservable struct {
	Token    token.Token
	Left     observer.ObservableNode[ast.Node]
	Operator string
	Right    observer.ObservableNode[ast.Node]
}

func (ie *InfixExpressionObservable) expressionNode()      {}
func (ie *InfixExpressionObservable) TokenLiteral() string { return ie.Token.Literal }

func (ie *InfixExpressionObservable) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.Node.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.Node.String())
	out.WriteString(")")

	return out.String()
}
