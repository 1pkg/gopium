package fmtio

import (
	"bytes"
	"io"
	"os"
)

// list of supported extensions
const (
	GOPIUM = "gopium"
	GO     = "go"
	JSON   = "json"
	XML    = "xml"
	CSV    = "csv"
	MD     = "md"
	HTML   = "html"
)

// stdout defines tiny wrapper for
// os stdout stream that couldn't be closed
type stdout struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Write just reuses os stdout stream write
func (stdout) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// Close just does nothing
func (stdout) Close() error {
	return nil
}

// Buffer defines buffer creator helper
func Buffer() io.ReadWriter {
	return &bytes.Buffer{}
}
