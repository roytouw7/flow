package object

import "Flow/src/utility/observer"

type ObservableObject[T Object] struct {
	Object T
	*observer.Signal
}

func WrapObjectWithObservable[T Object](object T, observer *observer.Signal) ObservableObject[T] {
	return ObservableObject[T]{
		Object: object,
		Signal: observer,
	}
}
