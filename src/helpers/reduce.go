package helpers

// Reduce reduces slice by function of T to U
// seed U will be used as the seed of the output
func Reduce[T, U any](slice []T, reducer func(result U, value T) U, seed U) U {
	result := seed

	for i := 0; i < len(slice); i++ {
		result = reducer(result, slice[i])
	}

	return result
}

// MapReduce reduces slice by function from T to V
// mapper maps T to U, the intermediate value
// reducer reduces multiple U to a single V
// seed V will be used as the seed of the output
func MapReduce[T, U, V any](slice []T, mapper func(in T) U, reducer func(result V, intermediate U) V, seed V) V {
	result := seed

	for i := 0; i < len(slice); i++ {
		result = reducer(result, mapper(slice[i]))
	}

	return result
}