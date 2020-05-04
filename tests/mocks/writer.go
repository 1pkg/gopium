package mocks

import (
	"bytes"
	"io"
	"sync"
)

// wc defines mock fmtio writer closer implementation
type wc struct {
	buf  *bytes.Buffer
	werr error
	cerr error
}

// Write mock implementation
func (wc wc) Write(p []byte) (n int, err error) {
	// in case we have error
	// return it back
	if wc.werr != nil {
		return 0, wc.werr
	}
	// otherwise use buf impl
	return wc.buf.Write(p)
}

// Close mock implementation
func (wc wc) Close() error {
	return wc.cerr
}

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
	return wc{buf: &buf, werr: w.Werr, cerr: w.Cerr}, w.Err
}
