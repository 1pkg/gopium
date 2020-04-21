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
type Writer func(id, loc string) (io.Writer, error)

// Stdout defines Writer implementation
// which only returns os stdout all the time
func Stdout(string, string) (io.Writer, error) {
	return os.Stdout, nil
}

// FileJson defines Writer implementation
// which creates json file
func FileJson(id, loc string) (io.Writer, error) {
	return file(id, loc, "json")
}

// FileXml defines Writer implementation
// which creates xml file
func FileXml(id, loc string) (io.Writer, error) {
	return file(id, loc, "xml")
}

// FileCsv defines Writer implementation
// which creates csv file
func FileCsv(id, loc string) (io.Writer, error) {
	return file(id, loc, "csv")
}

// FileGo defines Writer implementation
// which creates go file
func FileGo(id, loc string) (io.Writer, error) {
	return file(id, loc, "go")
}

// FileGopium defines Writer implementation
// which creates gopium file
func FileGopium(id, loc string) (io.Writer, error) {
	return file(id, loc, "gopium")
}

// file defines Writer helper
// which creates underlying file
// by name, path and ext
func file(name, path, ext string) (io.Writer, error) {
	bname := filepath.Base(name)
	bname = strings.Split(bname, ".")[0]
	path = filepath.Dir(path)
	return os.Create(fmt.Sprintf("%s/%s.%s", path, bname, ext))
}
