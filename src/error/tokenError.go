package cerr

import (
	"fmt"

	"Flow/src/token"
)

type tokenError struct {
	*baseError
	context *token.Token
}

func (t *tokenError) Error() string {
	fContext := fmt.Sprintf("%d:%d", t.context.Line, t.context.Pos)
	return fmt.Sprintf("%s: %s", fContext, t.err)
}
