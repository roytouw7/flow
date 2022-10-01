package cerr

import (
	"fmt"

	"Flow/src/token"
)

type compileErrorInterface interface {
	Error() string
	destruct() (*string, *token.Token)
	setError(err string)
}

type compileError struct {
	err         string
	fileContext *token.Token
}

func (p *compileError) Error() string {
	fContext := fmt.Sprintf("%d:%d", p.fileContext.Line, p.fileContext.Pos)
	return fmt.Sprintf("%s: %s", fContext, p.err)
}

func (p *compileError) destruct() (*string, *token.Token) {
	return &p.err, p.fileContext
}

func (p *compileError) setError(err string) {
	p.err = err
}

func Wrap[E compileErrorInterface](err E, context ...string) E {
	if len(context) == 1 {
		oldErr, _ := err.destruct()
		newErr := fmt.Sprintf("%s: %s", context[0], *oldErr)
		err.setError(newErr)

		return err
	}

	return Wrap(Wrap(err, context[len(context)-1]), context[:len(context)-1]...)
}
