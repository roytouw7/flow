package iterator

import (
	"fmt"
	"unicode"
)

type FileIterator interface {
	Next() (string, *MetaData, error)
	Peek() (string, error)
	PeekN(n int) (string, error)
	HasNext() bool
}

type fileIterator struct {
	source            string
	pos, line, relPos int
}

//todo refactor metadata out iterator, will be used in multiple places apart of iterator

type MetaData struct {
	pos, relPos, line int
}

func New(sourceFile string) FileIterator {
	return &fileIterator{
		source: sourceFile,
		pos:    0,
		relPos: 1,
		line:   1,
	}
}

func (iterator *fileIterator) Next() (string, *MetaData, error) {
	next, err := iterator.getNextValidCharacter()
	if err != nil {
		return "", nil, err
	}

	meta := iterator.getMetaData()
	iterator.incrementPosition()

	return next, meta, nil
}

func (iterator *fileIterator) getMetaData() *MetaData {
	return &MetaData{
		iterator.pos,
		iterator.relPos,
		iterator.line,
	}
}

func (iterator *fileIterator) getNextValidCharacter() (string, error) {
	for ok, err := isValidCharacter(iterator.currentChar()); !ok; ok, err = isValidCharacter(iterator.currentChar()) {
		if err != nil {
			return "", err
		}
		iterator.incrementPosition()
	}

	return iterator.currentChar(), nil
}

func (iterator *fileIterator) currentChar() string {
	return iterator.source[iterator.pos : iterator.pos+1]
}

func isValidCharacter(c string) (bool, error) {
	runes := []rune(c)
	if len(runes) > 1 {
		return false, fmt.Errorf("can not convert string %s to rune, too long", c)
	}

	if unicode.IsSpace(runes[0]) {
		return false, nil
	}

	return true, nil
}

func (iterator *fileIterator) incrementPosition() {
	iterator.relPos++

	if []rune(iterator.currentChar())[0] == '\n' {
		iterator.line++
		iterator.relPos = 1
	}

	iterator.pos++
}

func (iterator *fileIterator) HasNext() bool {
	return iterator.pos < len(iterator.source)
}

func (iterator *fileIterator) Peek() (string, error) {
	return iterator.peek(1)
}

func (iterator *fileIterator) PeekN(n int) (string, error) {
	return iterator.peek(n)
}

func (iterator *fileIterator) peek(n int) (string, error) {
	if iterator.pos+n > len(iterator.source) {
		return "", fmt.Errorf("peek %d out of bounds", n)
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
		return "", fmt.Errorf("peek %d out of bounds", n)
	}

	peekChar := iterator.source[offset : offset+1]
	if peekChar == "\r" {
		peekChar = "\n"
	}
	return peekChar, nil
}
