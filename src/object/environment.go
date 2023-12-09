package object

import "Flow/src/ast"

type Environment struct {
	store map[string]*ast.Expression
}

func NewEnvironment() *Environment {
	s := make(map[string]*ast.Expression)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (*ast.Expression, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val *ast.Expression) *ast.Expression {
	e.store[name] = val
	return val
}
