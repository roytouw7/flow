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

		if token == "" {
			test.T().Errorf("token is nil")
		}
		if err != nil {
			test.T().Error(err)
		}

		if token != expected {
			test.T().Errorf("expected token=%s got=%s", expected, token)
		}

		i++
	}

	if i != len(tt) {
		panic("not all test cases evaluated")
	}
}

func (test *Suite) TestMetaData() {
	data, err := os.ReadFile("test_assets/test_program.flow")
	if err != nil {
		panic(err)
	}
	iterator := New(string(data))

	tt := []struct {
		char              string
		pos, relPos, line int
	}{
		{"l", 0, 1, 1},
		{"e", 1, 2, 1},
		{"t", 2, 3, 1},
		{"f", 4, 5, 1},
		{"i", 5, 6, 1},
		{"v", 6, 7, 1},
		{"e", 7, 8, 1},
		{"7", 8, 9, 1},
		{"=", 10, 11, 1},
		{"9", 12, 13, 1},
		{";", 13, 14, 1},
		{"x", 16, 1, 2},
		{"+", 17, 2, 2},
		{"+", 18, 3, 2},
		{"c", 21, 1, 3},
		{"=", 23, 3, 3},
		{"=", 24, 4, 3},
		{"2", 26, 6, 3},
		{";", 27, 7, 3},
	}

	for _, t := range tt {
		_, metaData, err := iterator.Next()
		if err != nil {
			test.T().Error(err)
		}

		if metaData.pos != t.pos {
			test.T().Errorf("expected position to be %d, got %d", t.pos, metaData.pos)
		}
		if metaData.relPos != t.relPos {
			test.T().Errorf("expected relative position to be %d, got %d", t.relPos, metaData.relPos)
		}
		if metaData.line != t.line {
			test.T().Errorf("expected relative position to be %d, got %d", t.line, metaData.line)
		}
	}
}
