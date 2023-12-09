package object

import (
	"fmt"

	"Flow/src/utility/observer"
)

const OBSERVABLE_OBJ = "OBSERVABLE"

// todo value can either be a primitive or a composition of identifiers and optional primitives
// let a = 5;
// let b = a + 7;
// let c = (a + 13) * b / 2; => let c = composition(observables: [a, b], primitive) so basically we need lazy evaluation?

type Observable struct {
	Value Object // todo we no longer can save the value atomic, we need to be able to destruct it to apply changes from observable updates... maybe add a complex value representation or something?
	observer.BaseObservable[Object]
}

func NewObservable(obj Object) *Observable {
	return &Observable{
		Value: obj,
		BaseObservable: observer.BaseObservable[Object]{
			Observers: make([]observer.Observer[Object], 0),
		},
	}
}

func (o *Observable) Type() ObjectType {
	return OBSERVABLE_OBJ
}

func (o *Observable) Inspect() string {
	return fmt.Sprintf("%s observable", o.Value.Inspect())
}

func (o *Observable) Notify(change Object) {
	o.Value = change // todo this does not work because the underlying value could be a + 7, if a updates we need to add 7 again to a and set it as the value
	o.NotifyAll(o.Value)
}
