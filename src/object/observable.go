package object

import (
	"Flow/src/ast"
	"Flow/src/utility/observer"
)

const OBSERVABLE_OBJ = "OBSERVABLE"

type Observable struct {
	Value *ast.Expression // todo we no longer can save the value atomic, we need to be able to destruct it to apply changes from observable updates... maybe add a complex value representation or something?
	observer.BaseObservable[*ast.Expression]
}

func NewObservable(obj *ast.Expression) *Observable {
	return &Observable{
		Value: obj,
		BaseObservable: observer.BaseObservable[*ast.Expression]{
			Observers: make([]observer.Observer[*ast.Expression], 0),
		},
	}
}

func (o *Observable) Type() ObjectType {
	return OBSERVABLE_OBJ
}

func (o *Observable) Inspect() string {
	//return fmt.Sprintf("%s observable", o.Value.Inspect())
	return "" // todo
}

func (o *Observable) Notify(change *ast.Expression) {
	o.Value = change
	o.NotifyAll(o.Value)
}
