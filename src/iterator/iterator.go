package iterator

import (
	"fmt"
)

// StringIterator is the interface that wraps logic around iterating over strings
// Next returns the next smallest piece of string as a rune and increments its position
// Peek returns the next smallest piece of string without incrementing position
// HasNext safely checks if the next character can be retrieved or peeked
type StringIterator interface {
	Next() (rune, *MetaData, error)
	Peek() (rune, error)
	PeekN(n int) (rune, error)
	HasNext() bool
	HasNextN(n int) bool
}

type stringIterator struct {
	source            string
	pos, line, relPos int
}

//todo refactor metadata out iterator, will be used in multiple places apart of iterator

type MetaData struct {
	Pos, RelPos, Line int
}

func New(sourceFile string) StringIterator {
	return &stringIterator{
		source: sourceFile,
		pos:    0,
		relPos: 1,
		line:   1,
	}
}

func (iterator *stringIterator) Next() (rune, *MetaData, error) {
	next, err := iterator.getNextValidCharacter()
	if err != nil {
		return 0, nil, err
	}

	meta := iterator.getMetaData()
	iterator.incrementPosition()

	return next, meta, nil
}

func (iterator *stringIterator) getMetaData() *MetaData {
	return &MetaData{
		iterator.pos,
		iterator.relPos,
		iterator.line,
	}
}

func (iterator *stringIterator) getNextValidCharacter() (rune, error) {
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

func isValidCharacter(ch rune) (bool, error) {
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

func (iterator *stringIterator) Peek() (rune, error) {
	return iterator.peek(1)
}

func (iterator *stringIterator) PeekN(n int) (rune, error) {
	return iterator.peek(n)
}

func (iterator *stringIterator) peek(n int) (rune, error) {
	if iterator.pos+n > len(iterator.source) {
		return 0, fmt.Errorf("peek %d out of bounds", n)
	}

	// count newline characters "\r\n" as single increment
	offset := iterator.pos
	for i := 0; i < n; {
		if offset+1 > len(iterator.source) {
			return 0, fmt.Errorf("peek %d out of bounds", n)
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
