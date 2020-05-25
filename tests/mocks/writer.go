package mocks

import (
	"io"
	"os"
	"strings"
	"sync"

	"1pkg/gopium"
)

// Writer defines mock writer implementation
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
	// replace os path separators with underscores
	loc = strings.Replace(loc, gopium.Root(), "", 1)
	loc = strings.ReplaceAll(loc, string(os.PathSeparator), "_")
	loc = strings.Trim(loc, "_")
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

// Catwriter defines cat writer implementation
type Catwriter struct {
	Err error
}

// Catwriter mock implementation
func (cw Catwriter) Catwriter(w gopium.Writer, rcat string) (gopium.Writer, error) {
	return w, cw.Err
}
