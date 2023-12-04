package cerr

import (
	"fmt"
)

type baseErrorInterface interface {
	withoutContext() *string
	setError(err string)
}

type baseError struct {
	err         string
}

func (p *baseError) withoutContext() *string {
	return &p.err
}

func (p *baseError) setError(err string) {
	p.err = err
}

// Wrap wraps 1...n contexts to an error
func Wrap[E baseErrorInterface](err E, context ...string) E {
	if len(context) == 1 {
		oldErr := err.withoutContext()
		newErr := fmt.Sprintf("%s: %s", context[0], *oldErr)
		err.setError(newErr)

		return err
	}

	return Wrap(Wrap(err, context[len(context)-1]), context[:len(context)-1]...)
}
