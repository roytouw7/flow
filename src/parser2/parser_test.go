package parser2

import (
	"Flow/src/lexer"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestLetStatements() {
	data, err := os.ReadFile("test_assets/let_statements.flow")
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(data))
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(test.T(), p)
	if program == nil {
		test.Fail("ParseProgram() returned nil")
		return
	}

	checkProgramLines(test.T(), program, 3)

	tests := []struct {
		expectedIdentifier string
		expectedValue      int
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(test.T(), stmt, tt.expectedIdentifier, tt.expectedValue) {
			return
		}
	}
}
