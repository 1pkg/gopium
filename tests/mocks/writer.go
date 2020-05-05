package mocks

import (
	"bytes"
	"io"
	"sync"
)

// Writer defines mock fmtio writer implementation
type Writer struct {
	Buffers sync.Map
	Err     error
	Werr    error
	Cerr    error
}

// Writer mock implementation
func (w *Writer) Writer(id string, loc string) (io.WriteCloser, error) {
	// prepare shared buf
	var buf bytes.Buffer
	// write it to store
	// and to mock write closer
	w.Buffers.Store(id, &buf)
	return RWC{buf: &buf, Werr: w.Werr, Cerr: w.Cerr}, w.Err
}
