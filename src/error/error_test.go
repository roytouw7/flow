package cerr

import (
	"testing"

	"Flow/src/token"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestMissingParseFnError() {
	tok := token.Token{
		Line:    1,
		Pos:     7,
		Literal: token.ASTERISK,
		Type:    token.ASTERISK,
	}

	err := MissingParseFnError(&tok, Prefix)

	const expected = "1:7: no prefix parse function found for token \"*\""

	if err.Error() != expected {
		test.T().Errorf("MissingParseFnError constructed incorrect, expected=%s, got=%s", expected, err.Error())
	}
}

func (test *Suite) TestPeekOutOfBoundsError() {
	const (
		source = "source"
		line   = 1
		pos    = 7
		peek   = 14
	)

	err := PeekOutOfBoundsError(source, line, pos, peek)

	const expected = "source:1:7: peek out of bounds, trying to peek 14"

	if err.Error() != expected {
		test.T().Errorf("PeekOutOfBoundsError constructed incorrect, expected %q, got %q", expected, err.Error())
	}
}

func (test *Suite) TestWrapping() {
	tok := token.Token{
		Line:    1,
		Pos:     7,
		Literal: token.ASTERISK,
		Type:    token.ASTERISK,
	}

	err := Wrap(MissingParseFnError(&tok, Infix), "TestWrapping")

	var expected = "1:7: TestWrapping: no infix parse function found for token \"*\""

	if err.Error() != expected {
		test.T().Errorf("Wrapping error not working correct, expected=%s, got=%s", expected, err.Error())
	}

	err = Wrap(err, "LevelTwo")

	expected = "1:7: LevelTwo: TestWrapping: no infix parse function found for token \"*\""

	if err.Error() != expected {
		test.T().Errorf("Wrapping error second time not working correct, expected=%s, got=%s", expected, err.Error())
	}
}

func (test *Suite) TestWrappingMultipleContexts() {
	tok := token.Token{
		Line:    1,
		Pos:     7,
		Literal: token.ASTERISK,
		Type:    token.ASTERISK,
	}

	err := Wrap(MissingParseFnError(&tok, Infix), "Context1", "Context2", "Context3")

	var expected = "1:7: Context1: Context2: Context3: no infix parse function found for token \"*\""

	if err.Error() != expected {
		test.T().Errorf("Wrapping error not working correct, expected=%s, got=%s", expected, err.Error())
	}
}