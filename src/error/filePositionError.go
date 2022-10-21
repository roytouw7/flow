package cerr

import (
	"fmt"

	"Flow/src/metadata"
)

type filePositionError struct {
	*baseError
	metaData *metadata.MetaData
}

func (f *filePositionError) Error() string {
	m :=f.metaData
	fContext := fmt.Sprintf("%s:%d:%d", m.Source, m.Line, m.Pos)
	return fmt.Sprintf("%s: %s", fContext, f.err)
}
