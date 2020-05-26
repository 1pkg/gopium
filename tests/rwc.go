package tests

import "bytes"

// RWC defines mock io reader writer closer implementation
type RWC struct {
	buf  bytes.Buffer
	Rerr error
	Werr error
	Cerr error
}

// Read mock implementation
func (rwc *RWC) Read(p []byte) (int, error) {
	// in case we have error
	// return it back
	if rwc.Rerr != nil {
		return 0, rwc.Rerr
	}
	// otherwise use buf impl
	return rwc.buf.Read(p)
}

// Write mock implementation
func (rwc *RWC) Write(p []byte) (n int, err error) {
	// in case we have error
	// return it back
	if rwc.Werr != nil {
		return 0, rwc.Werr
	}
	// otherwise use buf impl
	return rwc.buf.Write(p)
}

// Close mock implementation
func (rwc *RWC) Close() error {
	return rwc.Cerr
}
