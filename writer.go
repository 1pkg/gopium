package gopium

import "io"

// Writer defines abstraction for
// io witer generation from set of parametrs
type Writer func(string) (io.WriteCloser, error)

// Catwriter defines abstraction for
// for writer with flexible root category
type Catwriter func(Writer, string) (Writer, error)
