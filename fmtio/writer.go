package fmtio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Writer defines abstraction for
// io witer generation from set of parametrs
type Writer func(id string, loc string) (io.WriteCloser, error)

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
func Stdout(string, string) (io.WriteCloser, error) {
	return stdout{}, nil
}

// File defines writer helper
// which creates underlying file callback
// by name, path and capture ext
func File(ext string) Writer {
	return func(id string, loc string) (io.WriteCloser, error) {
		bname := filepath.Base(id)
		bname = strings.Split(bname, ".")[0]
		path := filepath.Dir(loc)
		return os.Create(fmt.Sprintf("%s/%s.%s", path, bname, ext))
	}
}
