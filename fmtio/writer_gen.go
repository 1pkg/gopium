package fmtio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// WriteGen defines abstraction for
// io witer generation from set of parametrs
type WriterGen func(id, loc string) (io.Writer, error)

// Stdout defines WriterGen implementation
// which only returns os stdout all the time
func Stdout(string, string) (io.Writer, error) {
	return os.Stdout, nil
}

// TempFile defines WriterGen implementation
// which creates json file
func FileJson(id, loc string) (io.Writer, error) {
	return file(id, loc, "json")
}

// TempFile defines WriterGen implementation
// which creates xml file
func FileXml(id, loc string) (io.Writer, error) {
	return file(id, loc, "xml")
}

// TempFile defines WriterGen implementation
// which creates csv file
func FileCsv(id, loc string) (io.Writer, error) {
	return file(id, loc, "csv")
}

// TempFile defines WriterGen implementation
// which creates go file
func FileGo(id, loc string) (io.Writer, error) {
	return file(id, loc, "go")
}

// TempFile defines WriterGen implementation
// which creates gopium file
func FileGopium(id, loc string) (io.Writer, error) {
	return file(id, loc, "gopium")
}

// file defines WriterGen helper
// which creates underlying file
// by name, path and ext
func file(name, path, ext string) (io.Writer, error) {
	bname := filepath.Base(name)
	bname = strings.Split(bname, ".")[0]
	path = filepath.Dir(path)
	return os.Create(fmt.Sprintf("%s/%s.%s", path, bname, ext))
}
