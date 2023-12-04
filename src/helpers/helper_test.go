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

func TestReduce(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := "ac"

	reducer := func(result string, in int) string {
		switch in {
		case 1:
			return result + "a"
		case 3:
			return result + "c"
		default:
			return result
		}
	}

	result := Reduce(input, reducer, "")

	assert.Equal(t, expected, result)
	assert.Equal(t, []int{1, 2, 3, 4}, input)
}

func TestMapReduce(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := "ac"

	mapper := func(in int) string {
		switch in {
		case 1:
			return "a"
		case 3:
			return "c"
		default:
			return ""
		}
	}

	reducer := func(result string, intermediate string) string {
		return result + intermediate
	}

	result := MapReduce(input, mapper, reducer, "")

	assert.Equal(t, expected, result)
	assert.Equal(t, []int{1, 2, 3, 4}, input)
}

func TestMapReduceWithIntermediateState(t *testing.T) {
	input := []struct {
		val int
	}{{val: 1}, {val: 2}, {val: 3}}
	expected := "abc"

	mapper := func(in struct{ val int }) int {
		return in.val
	}

	reducer := func(result string, intermediate int) string {
		switch intermediate {
		case 1:
			return result + "a"
		case 2:
			return result + "b"
		case 3:
			return result + "c"
		default:
			return result
		}
	}

	result := MapReduce(input, mapper, reducer, "")

	assert.Equal(t, expected, result)
}
