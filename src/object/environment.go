package object

import "Flow/src/ast"

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
	if ok {	// if identifier points to itself, try finding identifier of outer scope to prevent infinite self referencing get
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

func (e *Environment) Set(name string, val *ast.Expression) *ast.Expression {
	e.store[name] = val
	return val
}
