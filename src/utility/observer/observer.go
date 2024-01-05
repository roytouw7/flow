package observer

import (
	"Flow/src/ast"
	"Flow/src/object"
)

type ObservableObserver interface {
	Observer
	Observable
}

type Observer interface {
	Notify()
}

type Observable interface {
	Register(observer Observer)
}

type ConcreteObservable struct {
	observers []Observer
}

func (c *ConcreteObservable) Register(observer Observer) {
	c.observers = append(c.observers, observer)
}

func (c *ConcreteObservable) Notify() {
	for _, observer := range c.observers {
		observer.Notify()
	}
}

type ObservableNode[T ast.Node] struct {
	Node T
	*ConcreteObservable
}

type ObservableObject[T object.Object] struct {
	Object T
	*ConcreteObservable
}

func WrapNodeWithObservable[T ast.Node](node T, observer *ConcreteObservable) ObservableNode[T] {
	return ObservableNode[T]{
		Node:               node,
		ConcreteObservable: observer,
	}
}

func WrapObjectWithObservable[T object.Object](object T, observer *ConcreteObservable) ObservableObject[T] {
	return ObservableObject[T]{
		Object:             object,
		ConcreteObservable: observer,
	}
}
