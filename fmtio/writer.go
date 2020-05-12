package fmtio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Writer defines abstraction for
// io witer generation from set of parametrs
type Writer func(string) (io.WriteCloser, error)

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

// Stdout defines writer implementation
// which only returns os stdout all the time
func Stdout(string) (io.WriteCloser, error) {
	return stdout{}, nil
}

// File defines writer helper
// which creates underlying file callback
// with capture name and ext on provided loc
func File(name string, ext string) Writer {
	return func(loc string) (io.WriteCloser, error) {
		path := filepath.Dir(loc)
		return os.Create(fmt.Sprintf("%s/%s.%s", path, name, ext))
	}
}

// Files defines writer helper
// which creates underlying files callback
// with capture ext on provided loc
func Files(ext string) Writer {
	return func(loc string) (io.WriteCloser, error) {
		path := strings.Replace(loc, filepath.Ext(loc), fmt.Sprintf(".%s", ext), 1)
		return os.Create(path)
	}
}

// Buffer defines buffer creator helper
func Buffer() io.ReadWriter {
	return &bytes.Buffer{}
}
