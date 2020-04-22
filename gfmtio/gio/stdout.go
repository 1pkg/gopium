package gio

import "os"

// Writer defines tiny wrapper for
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
