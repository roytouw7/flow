package object

import (
	"bytes"
	"strings"

	"Flow/src/utility/slice"
)

const ARRAY_OBJ = "ARRAY"

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := slice.Map(a.Elements, func(el Object) string {
		return el.Inspect()
	})

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
