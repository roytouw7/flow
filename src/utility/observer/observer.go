package observer

type Observer[T any] interface {
	Notify(T)
}

type Observable[T any] interface {
	Register(observer Observer[T])
	NotifyAll(change T)
	Unregister(observer Observer[T])
}

type BaseObservable[T any] struct {
	observers []Observer[T]
}

func (o *BaseObservable[T]) Register(observer Observer[T]) {
	o.observers = append(o.observers, observer)
}

func (o *BaseObservable[T]) NotifyAll(change T) {
	for _, observer := range o.observers {
		observer.Notify(change)
	}
}

func (o *BaseObservable[T]) Unregister(removeObservable Observer[T]) {
	for i, observer := range o.observers {
		if observer == removeObservable {
			o.observers = append(o.observers[:i], o.observers[i+1:]...)
			break
		}
	}
}
