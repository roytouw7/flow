package slice

// Map maps function over slice without mutating input
func Map[T, U any](slice []T, fn func(i T) U) []U {
	output := make([]U, len(slice))

	for i := 0; i < len(slice); i++ {
		output[i] = fn(slice[i])
	}

	return output
}