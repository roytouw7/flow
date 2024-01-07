package object

import (
	"fmt"

	"Flow/src/ast"
	"Flow/src/utility/observer"
)

type Environment struct {
	store           map[string]observer.ObservableNode[ast.Node]
	nativeFunctions map[string]*NativeFunc
	outer           *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]observer.ObservableNode[ast.Node])
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get recursively tries finding value for name in environment and outer environments
func (e *Environment) Get(name string) (observer.ObservableNode[ast.Node], bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) mustGet(name string) observer.ObservableNode[ast.Node] {
	val, ok := e.Get(name)
	if !ok {
		panic(fmt.Sprintf("variable %q not found!", name))
	}

	return val
}

func (e *Environment) Set(name string, val observer.ObservableNode[ast.Node]) observer.ObservableNode[ast.Node] {
	e.store[name] = val
	return val
}

// SubstituteReferences substitutes identifiers with their value from the environment
// when a name is given only identifiers matching the name are substituted
func (e *Environment) SubstituteReferences(node observer.ObservableNode[ast.Node], name *string) observer.ObservableNode[ast.Node] {
	switch exp := node.Node.(type) {
	case *ast.IdentifierLiteral:
		if name == nil || *name == exp.Value {
			if _, ok := Builtins[exp.Value]; ok {
				return node
			}
			val, ok := e.Get(exp.Value)
			if !ok {
				panic(fmt.Sprintf("could not find identifier %s in closure or outer closures", exp.Value))
			}
			if e.outer != nil {
				return e.outer.SubstituteReferences(val, nil)
			}
			return val
		}
		return node
	case *ast.PrefixExpression:
		right := e.SubstituteReferences(observer.WrapNodeWithObservable[ast.Node](exp.Right, observer.New()), name)
		prefix := &ast.PrefixExpression{
			Token:    exp.Token,
			Operator: exp.Operator,
			Right:    right.Node.(ast.Expression),
		}
		return observer.WrapNodeWithObservable[ast.Node](prefix, observer.New())
	case *ast.InfixExpression:
		var left, right observer.ObservableNode[ast.Node]
		var existingObservable *observer.Signal
		if exp.Operator == "=" { // don't substitute left hand side of assignment expression
			if identifier, ok := exp.Left.(*ast.IdentifierLiteral); ok {
				wrapped := observer.WrapNodeWithObservable[ast.Node](exp.Right, observer.New())
				right = e.SubstituteReferences(wrapped, &identifier.Value)
				existingObservable = e.mustGet(identifier.Value).Signal
			} else if _, ok := exp.Left.(*ast.IndexExpression); ok {
				wrapped := observer.WrapNodeWithObservable[ast.Node](exp.Right, observer.New())
				right = e.SubstituteReferences(wrapped, nil)
			} else {
				panic(fmt.Sprintf("expected left hand side of assignment expression to be identifier literal, got=%T", exp.Left))
			}
			if existingObservable != nil {
				left = observer.WrapNodeWithObservable[ast.Node](exp.Left, existingObservable)
			} else {
				left = observer.WrapNodeWithObservable[ast.Node](exp.Left, observer.New())
			}
		} else {
			left = e.SubstituteReferences(observer.WrapNodeWithObservable[ast.Node](exp.Left, observer.New()), name)
			right = e.SubstituteReferences(observer.WrapNodeWithObservable[ast.Node](exp.Right, observer.New()), name)
		}
		infix := &InfixExpressionObservable{
			Token:    exp.Token,
			Left:     left,
			Operator: exp.Operator,
			Right:    right,
		}
		return observer.WrapNodeWithObservable[ast.Node](infix, observer.New())
	case *ast.SliceLiteral:
		left := e.SubstituteReferences(observer.WrapNodeWithObservable[ast.Node](exp.Left, observer.New()), name)
		sliceLiteral := &ast.SliceLiteral{
			Token: exp.Token,
			Left:  left.Node.(ast.Expression),
			Lower: exp.Lower,
			Upper: exp.Upper,
		}
		return observer.WrapNodeWithObservable[ast.Node](sliceLiteral, observer.New())
	case *ast.IndexExpression:
		left := e.SubstituteReferences(observer.WrapNodeWithObservable[ast.Node](exp.Left, observer.New()), name)
		indexExp := &ast.IndexExpression{
			Token: exp.Token,
			Left:  left.Node.(ast.Expression),
			Index: exp.Index,
		}
		return observer.WrapNodeWithObservable[ast.Node](indexExp, observer.New())
	default:
		return node
	}
}
