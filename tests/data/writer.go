package data

import (
	"io"

	"1pkg/gopium"
)

// Writer defines tests data writter implementation
// which reuses underlying locator
// but purifies location generation
type Writer struct {
	Writer gopium.CategoryWriter
}

// Generate writer implementation
func (w Writer) Generate(loc string) (io.WriteCloser, error) {
	// purify the loc then
	// just reuse underlying writer
	return w.Writer.Generate(purify(loc))
}

// Category writer implementation
func (w Writer) Category(cat string) error {
	return w.Writer.Category(cat)
}
