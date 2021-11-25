package main

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
}

func (test Suite) SetupTest() {}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestRandChar() {
	for i := 0; i < 1000; i++ {
		r := randChar()
		if !isValidIdentifierChar(r) {
			test.Failf("Invalid character %s", string(r))
		}
	}
}

func (test *Suite) TestGenerateIdentifier() {
	identifier := generateIdentifier(5, 10)
	if len(identifier) < 5 || len(identifier) > 10 {
		msg := fmt.Sprintf("Identifier %s fails in length: %d", identifier, len(identifier))
		test.Fail(msg)
	}

	for i, k := range identifier {
		if !isValidIdentifierChar(k) {
			test.Failf("Invalid character %s in identifier %s at index %d", string(k), identifier, i)
		}
	}
}

func isValidIdentifierChar(char rune) bool {
	if int(char) < 65 || int(char) > 90 { // a...z
		if int(char) < 97 || int(char) > 122 { // A...Z
			if int(char) != 95 { // _
				return false
			}
		}
	}

	return true
}
