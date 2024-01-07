package observer

import "Flow/src/ast"

type ObservableObserver interface {
	Observer
	Observable
}

type Observer interface {
	Notify(id *TraceId)
	notify(id TraceId)
	SetHandler(func(id TraceId))
}

type Observable interface {
	Register(s *Signal)
}

type ObservableNode[T ast.Node] struct {
	Node T
	*Signal
}

func WrapNodeWithObservable[T ast.Node](node T, observer *Signal) ObservableNode[T] {
	return ObservableNode[T]{
		Node:   node,
		Signal: observer,
	}
}

//type ConcreteObservable struct {
//	observers []Observer
//	handler   *func(id int)
//	id        int
//}

//func (c *ConcreteObservable) Register(observer Observer) {
//	c.observers = append(c.observers, observer)
//}
//
//func (c *ConcreteObservable) Notify(id int) {
//	if c.handler != nil {
//		(*c.handler)(id)
//	}
//
//	for _, observer := range c.observers {
//		if id != observer.GetId() {
//				observer.Notify(id)
//		}
//	}
//}
//
//func (c *ConcreteObservable) SetHandler(fn func(id int)) {
//	c.handler = &fn
//}
//
//func (c *ConcreteObservable) GetId() int {
//	return c.id
//}
//
//type ObservableNode[T ast.Node] struct {
//	Node T
//	*ConcreteObservable
//}
//
//func WrapNodeWithObservable[T ast.Node](node T, observer *ConcreteObservable) ObservableNode[T] {
//	return ObservableNode[T]{
//		Node:               node,
//		ConcreteObservable: observer,
//	}
//}

//func New() *ConcreteObservable {
//	newId := idCounter
//	idCounter += 1
//
//	return &ConcreteObservable{
//		observers: make([]Observer, 0),
//		id:        newId,
//	}
//}
