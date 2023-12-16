package object

type NativeFunction func(args ...Object) Object

const NATIVE_FN_OBJ = "NATIVE_FN"

type NativeFunc struct {
	Fn NativeFunction
}

func (f *NativeFunc) Type() ObjectType {
	return NATIVE_FN_OBJ
}

func (f *NativeFunc) Inspect() string {
	return "native function"
}
