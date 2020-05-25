package gopium

import "io"

// Writer defines abstraction for
// io witers generation
type Writer interface {
	Generate(string) (io.WriteCloser, error)
}

// CategoryWriter defines abstraction for
// io witers generation with flexible category
type CategoryWriter interface {
	Category(string) error
	Writer
}
