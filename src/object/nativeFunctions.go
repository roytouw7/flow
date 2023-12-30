package object

import (
	"fmt"
)

// todo might not be the right package to place this

var Builtins = map[string]*NativeFunc{
	"len": {
		Fn: func(args ...Object) Object {
			return flowLen(args...)
		},
	},
	"print": {
		Fn: print,
	},
}

func flowLen(args ...Object) Object {
	if len(args) != 1 {
		return NewEvalErrorObject(fmt.Sprintf("expected 1 argument for len got=%d", len(args)))
	}
	switch arg := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(arg.Value))}
	case *Array:
		return &Integer{Value: int64(len(arg.Elements))}
	default:
		return NewEvalErrorObject(fmt.Sprintf("argument to \"len\" not supported, got=%T", args[0]))
	}
}

func print(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}
