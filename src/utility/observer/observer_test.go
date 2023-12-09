package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type concreteObserver struct {
	value int
}

func (c *concreteObserver) Notify(change int) {
	c.value = change
}

type concreteObservable struct {
	value int
	BaseObservable[int]
}

func (o *concreteObservable) update(change int) {
	o.value = change
	o.NotifyAll(change)
}

func (test *Suite) TestObserver() {
	observers := []*concreteObserver{
		{},
		{},
		{},
		{},
	}

	observable := concreteObservable{}

	for _, observer := range observers {
		observable.Register(observer)
	}

	observable.update(1)

	for _, observer := range observers {
		assert.Equal(test.T(), observable.value, observer.value)
	}

	observable.update(7)

	var counter int
	for _, observer := range observers {
		assert.Equal(test.T(), observable.value, observer.value)
		counter++
	}

	assert.Equal(test.T(), 4, counter)

	observable.Unregister(observers[0])

	assert.Equal(test.T(), len(observable.Observers), 3)

	observable.update(14)

	assert.Equal(test.T(), observers[0].value, 7)
	assert.Equal(test.T(), observers[1].value, 14)
	assert.Equal(test.T(), observers[2].value, 14)
	assert.Equal(test.T(), observers[3].value, 14)


}
