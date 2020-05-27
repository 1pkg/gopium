package mocks

import (
	"io"
	"sync"
)

// Writer defines mock category writer implementation
type Writer struct {
	Gerr  error
	Cerr  error
	RWCs  map[string]*RWC
	mutex sync.Mutex
}

// Generate mock implementation
func (w *Writer) Generate(loc string) (io.WriteCloser, error) {
	// lock rwcs access
	// and init them if they
	// haven't inited before
	defer w.mutex.Unlock()
	w.mutex.Lock()
	if w.RWCs == nil {
		w.RWCs = make(map[string]*RWC)
	}
	// if loc is inside existed rwcs
	// just return found rwc back
	if rwc, ok := w.RWCs[loc]; ok {
		return rwc, w.Gerr
	}
	// otherwise create new rwc
	// store and return it back
	rwc := &RWC{}
	w.RWCs[loc] = rwc
	return rwc, w.Gerr
}

// Category mock implementation
func (w *Writer) Category(cat string) error {
	return w.Cerr
}
