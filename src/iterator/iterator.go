package iterator

import (
	"fmt"
)

type StringIterator interface {
	Next() (rune, *MetaData, error)
	Peek() (rune, error)
	PeekN(n int) (rune, error)
	HasNext() bool
}

type stringIterator struct {
	source            string
	pos, line, relPos int
}

//todo refactor metadata out iterator, will be used in multiple places apart of iterator

type MetaData struct {
	pos, relPos, line int
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
	return iterator.pos < len(iterator.source)
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
		if []rune(iterator.source[offset : offset+1])[0] != '\r' {
			i++
		}
		offset++
	}

	if offset+1 > len(iterator.source) {
		return 0, fmt.Errorf("peek %d out of bounds", n)
	}

	peekChar := []rune(iterator.source[offset : offset+1])[0]
	if peekChar == '\r' {
		peekChar = '\n'
	}
	return peekChar, nil
}
