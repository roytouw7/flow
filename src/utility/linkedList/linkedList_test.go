package linkedList

import (
	"testing"

	"Flow/src/utility/convert"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestLinkedList() {
	l := &LinkedList[int]{}
	l.Push(1)
	l.Push(2)
	l.Push(3)
	l.Push(4)

	tests := []struct {
		input *int
	}{
		{convert.NewInt(1)},
		{convert.NewInt(2)},
		{convert.NewInt(3)},
		{convert.NewInt(4)},
	}

	for i := 0; l.HasNext(); l = l.Next() {
		test.Equal(tests[i].input, l.Value)
		i++
	}
}
