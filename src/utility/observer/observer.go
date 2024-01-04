package observer

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

type concreteObservable struct {
	observers []Observer
}

func (c *concreteObservable) Register(observer Observer) {
	c.observers = append(c.observers, observer)
}

func (c *concreteObservable) Notify() {
	for _, observer := range c.observers {
		observer.Notify()
	}
}
