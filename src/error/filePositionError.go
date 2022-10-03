package cerr

import "fmt"

type filePositionError struct {
	*baseError
	line, pos int
	source    string
}

func (f *filePositionError) Error() string {
	fContext := fmt.Sprintf("%s:%d:%d", f.source, f.line, f.pos)
	return fmt.Sprintf("%s: %s", fContext, f.err)
}
