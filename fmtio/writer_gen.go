package fmtio

import (
	"fmt"
	"io"
	"os"
)

// WriteGen defines abstraction for
// io witer generation from set of parametrs
type WriterGen func(id, loc, tp string) (io.Writer, error)

// Stdout defines WriterGen implementation
// which only returns os stdout all the time
func Stdout(string, string, string) (io.Writer, error) {
	return os.Stdout, nil
}

// TempFile defines WriterGen implementation
// which uses underlying ioutil tempfile
func TempFile(id, loc, tp string) (io.Writer, error) {
	return os.Create(fmt.Sprintf("%s/%s.%s", loc, id, tp))
}
