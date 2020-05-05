package mocks

import (
	"bytes"
	"io"
	"sync"
)

// Writer defines mock fmtio writer implementation
type Writer struct {
	Buffers map[string]*bytes.Buffer
	mutex   sync.Mutex
	Err     error
	Werr    error
	Cerr    error
}

// Writer mock implementation
func (w *Writer) Writer(id string, loc string) (io.WriteCloser, error) {
	// lock buffers access
	// and init them
	defer w.mutex.Unlock()
	w.mutex.Lock()
	if w.Buffers == nil {
		w.Buffers = make(map[string]*bytes.Buffer)
	}
	// prepare shared buf
	var buf bytes.Buffer
	// write it to store
	// and to mock write closer
	w.Buffers[id] = &buf
	return RWC{buf: &buf, Werr: w.Werr, Cerr: w.Cerr}, w.Err
}
