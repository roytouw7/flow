package eval

import (
	"testing"

	"Flow/src/object"
	"Flow/src/parser"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (test *Suite) TestEvalIntegerExpression() {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-7", -7},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, 1, env)
		testIntegerObject(test.T(), evaluated, tt.expected)
	}
}

func (test *Suite) TestEvalBooleanExpression() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true == false", false},
		{"false == false", true},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, 1, env)
		testBooleanObject(test.T(), evaluated, tt.expected)
	}
}

func (test *Suite) TestIfElseExpressions() {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, 1, env)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(test.T(), evaluated, int64(integer))
		} else {
			testNullObject(test.T(), evaluated)
		}
	}
}

func (test *Suite) TestLetStatements() {
	tests := []struct {
		input    string
		expected int64
		stmts    int
	}{
		{"let a = 5; a;", 5, 2},
		{"let a = 5 * 5; a;", 25, 2},
		{"let a = 5; let b = a; b;", 5, 3},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15, 4},
		{"let a = 5; a = 10; a", 10, 3},
		{"let a = 10; let b = 7; a = a + b; a;", 17, 4},	// todo fails on a; a is stored as an identifier referencing expression a, which causes an infinite loop, nice...
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, tt.stmts, env)
		unwrapped := unwrapObservable(evaluated, env)

		testIntegerObject(test.T(), *unwrapped, tt.expected)
	}
}

func (test *Suite) TestErrorHandling() {
	tests := []struct {
		input    string
		expected string
		stmts    int
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN", 1},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN", 2},
		{"-true;", "1:1: unknown operator: -BOOLEAN", 1},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN", 1},
		{"5; true + false; 5;", "unknown operator: BOOLEAN + BOOLEAN", 3},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN", 1},
		{"foobar", "identifier not found: foobar", 1},
		{"7 = 9;", "can't assign to non-identifier type, got=*ast.IntegerLiteral", 1},
		{"a = 9;", "identifier not found: \"a\"", 1},
		{"let a = a;", "identifier not found: a", 1},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, tt.stmts, env)

		errObj, ok := evaluated.(*object.EvalError)
		if !ok {
			test.T().Errorf("no error object returned, got=%T (%+v)", evaluated, evaluated)
		}

		if errObj.Message != tt.expected {
			test.T().Errorf("wrong error message, expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func (test *Suite) TestBangOperator() {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!true", false},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!7", true},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, 1, env)
		testBooleanObject(test.T(), evaluated, tt.expected)
	}
}

func (test *Suite) TestReturnStatements() {
	tests := []struct {
		input    string
		expected int64
		stmts    int
	}{
		{"return 10;", 10, 1},
		{"return 10; 9;", 10, 2},
		{"return 2 * 5; 9;", 10, 2},
		{"9; return 2 * 5; 9;", 10, 3},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, tt.stmts, env)
		testIntegerObject(test.T(), evaluated, tt.expected)
	}

	program := parser.CreateProgramFromFile(test.T(), "./test_assets/return_statements.flow", 1)
	env := object.NewEnvironment()
	result := Eval(program, env)
	testIntegerObject(test.T(), result, 10)
}

func testEval(t *testing.T, input string, expectedStatements int, env *object.Environment) object.Object {
	p := parser.CreateProgram(t, input, expectedStatements)
	return Eval(p, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, expected=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t, expected=%t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object is not NULL, got=%T (%+v)", obj, obj)
		return false
	}

	return true
}
