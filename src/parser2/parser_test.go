package parser2

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

// todo add more complex assignment values

func (test *Suite) TestLetStatements() {
	program := createProgram(test.T(), "test_assets/let_statements.flow", 3)

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

// todo add more complex expression returns

func (test *Suite) TestReturnStatements() {
	program := createProgram(test.T(), "test_assets/return_statements.flow", 3)

	tests := []struct {
		expectedReturnValue interface{}
	}{
		{5},
		{10},
		{993322},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testReturnStatement(test.T(), stmt, tt.expectedReturnValue)
	}
}
