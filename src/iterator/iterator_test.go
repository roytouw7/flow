package iterator

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestNext() {
	data, err := os.ReadFile("test_assets/test_program.flow")
	if err != nil {
		panic(err)
	}
	iterator := New(string(data))

	tt := []string{"l", "e", "t", "f", "i", "v", "e", "7", "=", "9", ";", "x", "+", "+", "c", "=", "=", "2", ";"}

	i := 0
	for iterator.HasNext() {
		token, _, err := iterator.Next()
		expected := tt[i]
		i++

		if token == "" {
			test.T().Errorf("token is nil")
		}
		if err != nil {
			test.T().Error(err)
		}

		if token != expected {
			test.T().Errorf("expected token=%s got=%s", expected, token)
		}
	}

	if i != len(tt) {
		panic("not all test cases evaluated")
	}
}
