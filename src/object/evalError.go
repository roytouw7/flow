package object

import "fmt"

const (
	ERROR_OBJ = "ERROR"
)

type EvalError struct {
	Message string
}

func (e *EvalError) Type() ObjectType {
	return ERROR_OBJ
}

func (e *EvalError) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

func NewEvalErrorObject(format string, a ...interface{}) *EvalError {
	return &EvalError{Message: fmt.Sprintf(format, a...)}
}
