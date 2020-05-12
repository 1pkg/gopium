package mocks

import (
	"go/build"
	"io"
	"strings"
	"sync"
)

// Writer defines mock fmtio writer implementation
type Writer struct {
	Err   error
	RWCs  map[string]*RWC
	mutex sync.Mutex
}

// Writer mock implementation
func (w *Writer) Writer(loc string) (io.WriteCloser, error) {
	// lock rwcs access
	// and init them if they
	// haven't inited before
	defer w.mutex.Unlock()
	w.mutex.Lock()
	if w.RWCs == nil {
		w.RWCs = make(map[string]*RWC)
	}
	// remove abs part from loc
	loc = strings.Replace(loc, build.Default.GOPATH, "", 1)
	// if loc is inside existed rwcs
	// just return found rwc back
	if rwc, ok := w.RWCs[loc]; ok {
		return rwc, w.Err
	}
	// otherwise create new rwc
	// store and return it back
	rwc := &RWC{}
	w.RWCs[loc] = rwc
	return rwc, w.Err
}
