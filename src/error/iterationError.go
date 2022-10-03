package cerr

import "fmt"

type IterationError interface {
	error
	iterationError() // to discriminate IterationError from other Errors when required
	baseErrorInterface
}

type iterationError struct {
	*filePositionError
}

func (i *iterationError) iterationError() {}

func PeekOutOfBoundsError(source string, line, pos, peek int) IterationError {
	msg := fmt.Sprintf("peek out of bounds, trying to peek %d", peek)
	return newIterationError(msg, source, line, pos)
}

func newIterationError(msg string, source string, line, pos int) *iterationError {
	return &iterationError{
		&filePositionError{
			baseError: &baseError{msg},
			line:      line,
			pos:       pos,
			source:    source,
		},
	}
}
