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

func (e *Environment) Get(name string) (*ast.Expression, bool) {
	obj, ok := e.store[name]
	if ok { // if identifier points to itself, try finding identifier of outer scope to prevent infinite self referencing get
		if idx, isIdx := (*obj).(*ast.IdentifierLiteral); isIdx {
			if idx.Value == name {
				ok = false
			}
		}
	}

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
	//if _, isIdx := (*val).(*ast.IdentifierLiteral); isIdx {
	//	panic(fmt.Sprintf("self referencing idx, %q", name))
	//}
	e.store[name] = val
	return val
}

func (e *Environment) replaceIdentifier(node ast.Expression) ast.Expression {
	switch node := node.(type) {
	case *ast.IdentifierLiteral:
		val := *e.mustGet(node.Value)
		return e.replaceIdentifier(val)
	case *ast.PrefixExpression:
		right := e.replaceIdentifier(node.Right)
		node.Right = right
		return node
	case *ast.InfixExpression:
		left := e.replaceIdentifier(node.Left)
		right := e.replaceIdentifier(node.Right)
		node.Left = left
		node.Right = right
		return node
	default:
		return node
	}
}

// ReplaceWithOuterScopeValue first try getting it from the current scope, then pass it in again with the outer scope
func (e *Environment) ReplaceWithOuterScopeValue(node ast.Expression) ast.Expression {
	if e.outer != nil {
		switch node := node.(type) {
		case *ast.IdentifierLiteral:
			val, ok := e.Get(node.Value)
			if ok {
				return e.outer.ReplaceWithOuterScopeValue(*val) // when value found pass this
			}
			return node // else pass identifier and try outer scope
		case *ast.PrefixExpression:
			right := e.ReplaceWithOuterScopeValue(node.Right)
			node.Right = right
			return node
		case *ast.InfixExpression:
			left := e.ReplaceWithOuterScopeValue(node.Left)
			node.Left = left
			right := e.ReplaceWithOuterScopeValue(node.Right)
			node.Right = right
			return node
		default:
			return node
		}
	}

	return e.replaceIdentifier(node)
}

func (e *Environment) ReplaceWithOuterScopeValue2(node ast.Expression) ast.Expression {
	if e.outer != nil {
		switch node := node.(type) {
		case *ast.IdentifierLiteral:
			val, ok := e.Get(node.Value)
			if ok {
				return e.outer.ReplaceWithOuterScopeValue2(*val) // when value found pass this
			}
			return node // else pass identifier and try outer scope

		case *ast.PrefixExpression:
			right := e.ReplaceWithOuterScopeValue2(node.Right)
			// Create a new PrefixExpression with the updated right expression
			return &ast.PrefixExpression{
				Operator: node.Operator,
				Right:    right,
			}

		case *ast.InfixExpression:
			left := e.ReplaceWithOuterScopeValue2(node.Left)
			right := e.ReplaceWithOuterScopeValue2(node.Right)
			// Create a new InfixExpression with the updated left and right expressions
			return &ast.InfixExpression{
				Operator: node.Operator,
				Left:     left,
				Right:    right,
			}

		default:
			return node
		}
	}

	return e.replaceIdentifier(node)
}
