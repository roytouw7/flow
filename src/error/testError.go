package cerr

import (
	"fmt"

	"Flow/src/token"
)

type TestError interface {
	error
	baseErrorInterface
}

type testError struct {
	*baseError
	testCase  string
	testIndex int
}

func (t *testError) testError() {}

func (t *testError) Error() string {
	testContext := fmt.Sprintf("%s[%d]", t.testCase, t.testIndex)
	return fmt.Sprintf("TestCase %s: %s", testContext, t.err)
}

func TestUnexpectedError(test string, testIndex int, err baseErrorInterface) TestError {
	msg := fmt.Sprintf("unexpected error %q", err)
	return newTestError(test, testIndex, msg)
}

func TestUnexpectedValueFor(test string, testIndex int, identifier string, actual, expected interface{}) TestError {
	msg := fmt.Sprintf("expected value %q for identifier %q, got %q", expected, identifier, actual)
	return newTestError(test, testIndex, msg)
}

func TestUnexpectedTokenError(test string, testIndex int, tok *token.Token, expected token.Type) TestError {
	msg := UnexpectedTokenError(tok, expected)
	return newTestError(test, testIndex, msg.Error())
}

func newTestError(test string, testIndex int, msg string) TestError {
	return &testError{
		baseError: &baseError{
			err: msg,
		},
		testCase:  test,
		testIndex: testIndex,
	}
}
