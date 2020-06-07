package fmtio

import (
	"bytes"
	"io"
	"os"
)

// list of supported extensions
const (
	JSON = "json"
	XML  = "xml"
	CSV  = "csv"
	MD   = "md"
)

// stdout defines tiny wrapper for
// os stdout stream that couldn't be closed
type stdout struct{}

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
