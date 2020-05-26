package mocks

import (
	"io"
	"os"
	"strings"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/tests"
)

// Writer defines mock category writer implementation
type Writer struct {
	Gerr  error
	Cerr  error
	RWCs  map[string]*tests.RWC
	mutex sync.Mutex
}

// Writer mock implementation
func (w *Writer) Generate(loc string) (io.WriteCloser, error) {
	// lock rwcs access
	// and init them if they
	// haven't inited before
	defer w.mutex.Unlock()
	w.mutex.Lock()
	if w.RWCs == nil {
		w.RWCs = make(map[string]*tests.RWC)
	}
	// remove abs part from loc
	// replace os path separators with underscores
	loc = strings.Replace(loc, gopium.Root(), "", 1)
	loc = strings.ReplaceAll(loc, string(os.PathSeparator), "_")
	loc = strings.Trim(loc, "_")
	// if loc is inside existed rwcs
	// just return found rwc back
	if rwc, ok := w.RWCs[loc]; ok {
		return rwc, w.Gerr
	}
	// otherwise create new rwc
	// store and return it back
	rwc := &tests.RWC{}
	w.RWCs[loc] = rwc
	return rwc, w.Gerr
}

// Category mock implementation
func (w *Writer) Category(cat string) error {
	return w.Cerr
}
