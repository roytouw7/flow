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
