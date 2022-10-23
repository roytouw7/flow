package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	result := Map(input, func(i int) int { return i * 2 })

	assert.Equal(t, expected, result)
	assert.Equal(t, []int{1, 2, 3, 4}, input)
}

func TestMapTU(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []string{"a", "b", "c", "d"}
	result := Map(input, func(i int) string {
		switch i {
		case 1:
			return "a"
		case 2:
			return "b"
		case 3:
			return "c"
		case 4:
			return "d"
		default:
			panic("unexpected input")
		}
	})

	assert.Equal(t, expected, result)
	assert.Equal(t, []int{1, 2, 3, 4}, input)
}