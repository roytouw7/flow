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
	Observers []Observer[T]
}

func (o *BaseObservable[T]) Register(observer Observer[T]) {
	o.Observers = append(o.Observers, observer)
}

func (o *BaseObservable[T]) NotifyAll(change T) {
	for _, observer := range o.Observers {
		observer.Notify(change)
	}
}

func (o *BaseObservable[T]) Unregister(removeObservable Observer[T]) {
	for i, observer := range o.Observers {
		if observer == removeObservable {
			o.Observers = append(o.Observers[:i], o.Observers[i+1:]...)
			break
		}
	}
}
