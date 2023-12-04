package iterator

import (
	"fmt"

	cerr "Flow/src/error"
	"Flow/src/metadata"
)

// todo this error handling is quite a mess a lot of errors are never set, the only type of error the iterator should return is IterationError
// StringIterator is the interface that wraps logic around iterating over strings
// Next returns the next smallest piece of string as a rune and increments its position
// Peek returns the next smallest piece of string without incrementing position
// HasNext safely checks if the next character can be retrieved or peeked
type StringIterator interface {
	Next() (rune, *metadata.MetaData, *cerr.IterationError)
	Peek() (rune, *cerr.IterationError)
	PeekN(n int) (rune, *cerr.IterationError)
	HasNext() bool
	HasNextN(n int) bool
}

type stringIterator struct {
	source            string
	pos, line, relPos int
}

func New(sourceFile string) StringIterator {
	return &stringIterator{
		source: sourceFile,
		pos:    0,
		relPos: 1,
		line:   1,
	}
}

// Copy make shallow copy of iterator
func Copy(iterator StringIterator) (StringIterator, error) {
	if strIterator, ok := iterator.(*stringIterator); ok {
		iteratorCopy := &stringIterator{
			source: strIterator.source,
			pos:    strIterator.pos,
			line:   strIterator.line,
			relPos: strIterator.relPos,
		}
		return iteratorCopy, nil
	}
	return nil, fmt.Errorf("failed making copy of iterator, expected *iterator.stringIterator, got %T", iterator)
}

func (iterator *stringIterator) Next() (rune, *metadata.MetaData, *cerr.IterationError) {
	next, err := iterator.getNextValidCharacter()
	if err != nil {
		return 0, nil, err
	}

	meta := iterator.getMetaData()
	iterator.incrementPosition()

	return next, meta, nil
}

func (iterator *stringIterator) getMetaData() *metadata.MetaData {
	return &metadata.MetaData{
		"",	// todo we gonna fix the source?
		iterator.pos,
		iterator.relPos,
		iterator.line,
	}
}

func (iterator *stringIterator) getNextValidCharacter() (rune, *cerr.IterationError) {
	for ok, err := isValidCharacter(iterator.currentChar()); !ok; ok, err = isValidCharacter(iterator.currentChar()) {
		if err != nil {
			return 0, err
		}
		iterator.incrementPosition()
	}

	return iterator.currentChar(), nil
}

func (iterator *stringIterator) currentChar() rune {
	return []rune(iterator.source[iterator.pos : iterator.pos+1])[0]
}

// todo this creates no error yet...
func isValidCharacter(ch rune) (bool, *cerr.IterationError) {
	if ch == '\r' {
		return false, nil
	}

	return true, nil
}

func (iterator *stringIterator) incrementPosition() {
	switch iterator.currentChar() {
	case '\n':
		iterator.line++
		iterator.relPos = 1
		iterator.pos++
	case '\r':
		iterator.pos++
	default:
		iterator.pos++
		iterator.relPos++
	}
}

func (iterator *stringIterator) HasNext() bool {
	return iterator.hasNext(1)
}

func (iterator *stringIterator) HasNextN(n int) bool {
	return iterator.hasNext(n)
}

func (iterator *stringIterator) hasNext(n int) bool {
	return iterator.pos+n-1 < len(iterator.source)
}

func (iterator *stringIterator) Peek() (rune, *cerr.IterationError) {
	return iterator.peek(1)
}

func (iterator *stringIterator) PeekN(n int) (rune, *cerr.IterationError) {
	return iterator.peek(n)
}

func (iterator *stringIterator) peek(n int) (rune, *cerr.IterationError) {
	if iterator.pos+n > len(iterator.source) {
		err := cerr.PeekOutOfBoundsError(iterator.source, iterator.line, iterator.pos, n)
		return 0, &err
	}

	// count newline characters "\r\n" as single increment
	offset := iterator.pos
	for i := 0; i < n; {
		if offset+1 > len(iterator.source) {
			err := cerr.PeekOutOfBoundsError(iterator.source, iterator.line, iterator.pos, n)
			return 0, &err
		}
		if []rune(iterator.source[offset : offset+1])[0] != '\r' {
			i++
		}
		offset++
	}

	peekChar := []rune(iterator.source[offset-1 : offset])[0]
	if peekChar == '\r' {
		peekChar = '\n'
	}
	return peekChar, nil
}
