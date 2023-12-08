package slice

import "testing"

var (
	input40 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}
	input20 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	input10 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	input1  = []int{1}

	fn      = func(i int) int {
		return i * 2
	}

	result []int
)

func benchmarkMap(input []int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		result = Map(input, fn)
	}
}

func benchmarkImperative(input []int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		output := make([]int, len(input))

		for i, number := range input {
			output[i] = fn(number)
		}

		result = output
	}
}

func BenchmarkMap40(b *testing.B) { benchmarkMap(input40, b) }
func BenchmarkMap20(b *testing.B) { benchmarkMap(input20, b) }
func BenchmarkMap10(b *testing.B) { benchmarkMap(input10, b) }
func BenchmarkMap1(b *testing.B)  { benchmarkMap(input1, b) }

func BenchmarkMapImperative40(b *testing.B) { benchmarkImperative(input40, b) }
func BenchmarkMapImperative20(b *testing.B) { benchmarkImperative(input20, b) }
func BenchmarkMapImperative10(b *testing.B) { benchmarkImperative(input10, b) }
func BenchmarkMapImperative1(b *testing.B)  { benchmarkImperative(input1, b) }
