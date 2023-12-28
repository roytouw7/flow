package object

import (
	"fmt"

	"Flow/src/ast"
)

type Environment struct {
	store map[string]*ast.Expression
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]*ast.Expression)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get recursively tries finding value for name in environment and outer environments
func (e *Environment) Get(name string) (*ast.Expression, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) mustGet(name string) *ast.Expression {
	val, ok := e.Get(name)
	if !ok {
		panic(fmt.Sprintf("variable %q not found!", name))
	}

	return val
}

func (e *Environment) Set(name string, val *ast.Expression) *ast.Expression {
	e.store[name] = val
	return val
}

// SubstituteReferences substitutes identifiers with their value from the environment
// when a name is given only identifiers matching the name are substituted
func (e *Environment) SubstituteReferences(node ast.Expression, name *string) ast.Expression {
	switch node := node.(type) {
	case *ast.IdentifierLiteral:
		if name == nil || *name == node.Value {
			val, ok := e.Get(node.Value)
			if !ok {
				panic(fmt.Sprintf("could not find identifier %s in environment or outer environments", node.Value))
			}
			if e.outer != nil {
				return e.outer.SubstituteReferences(*val, nil)
			}
			return *val
		}
		return node
	case *ast.PrefixExpression:
		right := e.SubstituteReferences(node.Right, name)
		return &ast.PrefixExpression{
			Token:    node.Token,
			Operator: node.Operator,
			Right:    right,
		}
	case *ast.InfixExpression:
		var left, right ast.Expression
		if node.Operator == "=" { // don't substitute left hand side of assignment expression
			if identifier, ok := node.Left.(*ast.IdentifierLiteral); ok {
				right = e.SubstituteReferences(node.Right, &identifier.Value)
			} else {
				panic(fmt.Sprintf("expected left hand side of assignment expression to be identifier literal, got=%T", node.Left))
			}
			left = node.Left
		} else {
			left = e.SubstituteReferences(node.Left, name)
			right = e.SubstituteReferences(node.Right, name)
		}
		return &ast.InfixExpression{
			Token:    node.Token,
			Left:     left,
			Operator: node.Operator,
			Right:    right,
		}
	case *ast.SliceLiteral:
		left := e.SubstituteReferences(node.Left, name)
		return &ast.SliceLiteral{
			Token: node.Token,
			Left:  left,
			Lower: node.Lower,
			Upper: node.Upper,
		}
	default:
		return node
	}
}
