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
	iterator := prepareIterator("test_assets/test_program.flow")

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
	iterator := prepareIterator("test_assets/test_program.flow")

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

func (test *Suite) TestPeek() {
	iterator := prepareIterator("test_assets/test_program.flow")

	char, err := iterator.Peek()
	test.expectChar(char, err, "e")
	char, err = iterator.PeekN(5)
	test.expectChar(char, err, "i")
	char, err = iterator.PeekN(14)
	test.expectChar(char, err, "\n")
	char, err = iterator.PeekN(19)
	test.expectChar(char, err, "c")

	_, _, _ = iterator.Next()
	char, err = iterator.Peek()
	test.expectChar(char, err, "t")
	char, err = iterator.PeekN(5)
	test.expectChar(char, err, "v")

	// move enough to be at second line of input
	for i := 0; i < 10; i++ {
		_, _, _ = iterator.Next()
	}
	char, err = iterator.Peek()
	test.expectChar(char, err, "x")

	// test out of bound peek handling
	_, err = iterator.PeekN(11)
	if err != nil {
		test.T().Errorf("unexpected out of bounds peek")
	}
	_, err = iterator.PeekN(12)
	if err == nil {
		test.T().Errorf("expected out of bounds error on peek but received nil")
	}

	// move to last character
	for iterator.HasNext() {
		_, _, _ = iterator.Next()
	}
	_, err = iterator.Peek()
	if err == nil {
		test.T().Errorf("expected out of bounds error on peek but received nil")
	}
}

func prepareIterator(sourceFile string) FileIterator {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		panic(err)
	}
	iterator := New(string(data))

	return iterator
}

func (test *Suite) expectChar(char string, err error, expected string) bool {
	if err != nil {
		test.T().Error(err)
		return false
	}

	if char != expected {
		test.T().Errorf("expected peek char to be %s got %s", expected, char)
		return false
	}

	return true
}
