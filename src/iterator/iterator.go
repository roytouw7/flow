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
	}
}

func (i *fileIterator) Next() (string, *MetaData, error) {
	next, err := i.getNextValidCharacter()
	if err != nil {
		return "", nil, err
	}

	i.incrementPosition()

	return next, nil, nil
}

func (i *fileIterator) getNextValidCharacter() (string, error) {
	for ok, err := isValidCharacter(i.source[i.pos : i.pos+1]); !ok; ok, err = isValidCharacter(i.source[i.pos : i.pos+1]) {
		if err != nil {
			return "", err
		}
		i.incrementPosition()
	}

	return i.source[i.pos : i.pos+1], nil
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

func (i *fileIterator) incrementPosition() {
	i.pos += 1
}

func (i *fileIterator) HasNext() bool {
	return i.pos < len(i.source)
}

func (i *fileIterator) Peek() (string, error) {
	return "", nil
}

func (i *fileIterator) PeekN(n int) (string, error) {
	return "", nil
}
