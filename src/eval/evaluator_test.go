package eval

import (
	"fmt"
	"testing"

	"Flow/src/ast"
	"Flow/src/object"
	"Flow/src/parser"

	"github.com/stretchr/testify/assert"
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
		{"let a = 10; let b = 7; a = a + b; a;", 17, 4},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, tt.stmts, env)

		testIntegerObject(test.T(), evaluated, tt.expected)
	}
}

func (test *Suite) TestFunctionLiterals() {
	input := "(x) => { x + 2; };"

	env := object.NewEnvironment()
	evaluated := testEval(test.T(), input, 1, env)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		test.Failf("TestFunctionLiterals", "object is not FunctionLiteral, got=%T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		test.Failf("TestFunctionLiterals", "expected 1 parameters, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		test.Failf("TestFunctionLiterals", "parameter is not \"x\" got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		test.Failf("TestFunctionLiterals", "expected body to be=%q, got=%q", expectedBody, fn.Body.String())
	}
}

func (test *Suite) TestFunctionApplication() {
	tests := []struct {
		input    string
		expected interface{}
		stmts    int
	}{
		{"let identity = (x) => { x; }; identity(5);", 5, 2},
		{"let identity = (x) => { return x; }; identity(5);", 5, 2},
		{"let add = (x, y) => { return x + y; }; add(7, 9);", 16, 2},
		{"let add = (x, y) => { x + y; }; add(5 + 5, add(5, 5));", 20, 2},
		{"(x) => { x; }(5);", 5, 1},
		{"let identity = (x) => { return x; }; identity(\"test\");", "test", 2},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, tt.stmts, env)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(test.T(), evaluated, int64(expected))
		case string:
			testStringObject(test.T(), evaluated, expected)
		}
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
		{"foobar", "Eval: failed substituting references for *ast.IdentifierLiteral \"foobar\", could not find identifier foobar in closure or outer closures", 1},
		{"7 = 9;", "Eval: failed substituting references for *ast.InfixExpression \"(7 = 9)\", expected left hand side of assignment expression to be identifier literal, got=*ast.IntegerLiteral", 1},
		{"a = 9;", "identifier not found: \"a\"", 1},
		{"let a = a;", "Eval: failed substituting references for *ast.IdentifierLiteral \"a\", could not find identifier a in closure or outer closures", 1},
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

func (test *Suite) TestStringLiteral() {
	p := parser.CreateProgram(test.T(), "\"foo ${1 + 7 * 9} bar\";", 1)
	env := object.NewEnvironment()
	result := Eval(p, env)
	testStringObject(test.T(), result, "foo 64 bar")
}

func (test *Suite) TestNativeFunctions() {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to \"len\" not supported, got=*object.Integer"},
		{`len("one", "two")`, "expected 1 argument for len got=2"},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, 1, env)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(test.T(), evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.EvalError)
			if !ok {
				test.T().Errorf("Object is not error, got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				test.T().Errorf("wrong error message, expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func (test *Suite) TestArrayLiteral() {
	input := "[1, 2 * 2, 3 + 3];"

	p := parser.CreateProgram(test.T(), input, 1)
	env := object.NewEnvironment()
	evaluated := Eval(p, env)

	array, ok := evaluated.(*object.Array)
	if !ok {
		test.T().Errorf("Object is not array, got=%T (%+v)", evaluated, evaluated)
	}

	if len(array.Elements) != 3 {
		test.T().Errorf("Expected array to have length 3, got=%d", len(array.Elements))
	}

	testIntegerObject(test.T(), array.Elements[0], 1)
	testIntegerObject(test.T(), array.Elements[1], 4)
	testIntegerObject(test.T(), array.Elements[2], 6)
}

func (test *Suite) TestArrayIndexing() {
	tests := []struct {
		input    string
		expected interface{}
		stmts    int
	}{
		{"[1, 2, 3][0]", 1, 1},
		{"[1, 2, 3][1]", 2, 1},
		{"[1, 2, 3][2]", 3, 1},
		{"let i = 0; [1][i]", 1, 2},
		{"[1, 2, 3][1 + 1]", 3, 1},
		{"let myArray = [1, 2, 3]; myArray[2];", 3, 2},
		{"let myArray = [1, 2, 3]; let i = myArray[0] + myArray[1] + myArray[2]; i;", 6, 3},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];", 2, 3},
		{"[1, 2, 3][3]", nil, 1},  // todo might want to return error here...?
		{"[1, 2, 3][-1]", nil, 1}, // todo might want to return 3 here...?
		{"[1, 2, 3][:];", []int{1, 2, 3}, 1},
		{"[1, 2, 3][1:];", []int{2, 3}, 1},
		{"[1, 2, 3][:2];", []int{1, 2}, 1},
		{"[1, 2, 3][1:2];", []int{2}, 1},
		{"let i = 0; let j = 2; [1, 2, 3][i:j];", []int{1, 2}, 3},
		{"let myArray = [1, 2, 3]; let lower = 1 + 0; let upper = 1 + 1; myArray[lower:upper];", []int{2}, 4},
		{"[1,2,3][1] = 7;", nil, 1},
		{"let arr = [1,2,3]; arr[1] = 7;", nil, 2},
		{"let arr = [1, 2, 3]; arr[2] = 7; arr[2];", 7, 3},
		{"let arr = [1, 2, 3]; let b = arr[:]; arr[0] = 7; b[0];", 1, 4},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		evaluated := testEval(test.T(), tt.input, tt.stmts, env)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(test.T(), evaluated, int64(expected))
		case []int:
			array := evaluated.(*object.Array)
			for i := 0; i < len(expected); i++ {
				testIntegerObject(test.T(), array.Elements[i], int64(expected[i]))
			}
		default:
			testNullObject(test.T(), evaluated)
		}
	}
}

func (test *Suite) TestEvaluatingHigherOrderFunctions() {
	program := parser.CreateProgramFromFile(test.T(), "./test_assets/higher_order_functions.flow", 2)
	env := object.NewEnvironment()
	evaluated := Eval(program, env)
	testIntegerObject(test.T(), evaluated, 0)
}

func (test *Suite) TestEvaluatingSum() {
	program := parser.CreateProgramFromFile(test.T(), "./test_assets/sum.flow", 2)
	env := object.NewEnvironment()
	evaluated := Eval(program, env)
	testIntegerObject(test.T(), evaluated, 6)
}

func (test *Suite) TestStringLiteralEvaluation() {
	// Set environment with 2 closures having n = (n - 1) - 1 with outer n 2
	env := object.NewEnvironment()
	env.Set("n", createInteger(2))
	env = object.NewEnclosedEnvironment(env)
	infix1 := createN1Infix("n", "-", *createInteger(1))
	env.Set("n", &infix1)
	env = object.NewEnclosedEnvironment(env)
	infix2 := createN1Infix("n", "-", *createInteger(1))
	env.Set("n", &infix2)

	p := parser.CreateProgram(test.T(), "\"n:${n}   n>0:${n>0}\"", 1)
	statement, ok := p.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected statement to be *ast.ExpressioNStatement, got=%T", statement))
	}
	stringLiteral, ok := statement.Expression.(*ast.StringLiteral)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected statement.Expression to be *ast.StringLiteral, got=%T", statement.Expression))
	}
	result := Eval(stringLiteral, env)
	stringResult, ok := result.(*object.String)
	if !ok {
		assert.Fail(test.T(), fmt.Sprintf("expected result to be *object.String, got=%T", result))
	}
	assert.Equal(test.T(), "n:0   n>0:false", stringResult.Value)
}
