package eval

import (
	"fmt"

	"Flow/src/object"
)

var builtins = map[string]*object.NativeFunc{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			return flowLen(args...)
		},
	},
	"print": {
		Fn: print,
	},
}

func flowLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newEvalErrorObject(fmt.Sprintf("expected 1 argument for len got=%d", len(args)))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newEvalErrorObject(fmt.Sprintf("argument to \"len\" not supported, got=%T", args[0]))
	}
}

func print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}
