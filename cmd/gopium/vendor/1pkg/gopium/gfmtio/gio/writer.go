package gio

import "io"

// Writer defines abstraction for
// io witer generation from set of parametrs
type Writer func(id, loc string) (io.WriteCloser, error)

// Stdout defines Writer implementation
// which only returns os stdout all the time
func Stdout(string, string) (io.WriteCloser, error) {
	return stdout{}, nil
}

// FileJson defines Writer implementation
// which creates json file
func FileJson(id, loc string) (io.WriteCloser, error) {
	return file(id, loc, "json")
}

// FileXml defines Writer implementation
// which creates xml file
func FileXml(id, loc string) (io.WriteCloser, error) {
	return file(id, loc, "xml")
}

// FileCsv defines Writer implementation
// which creates csv file
func FileCsv(id, loc string) (io.WriteCloser, error) {
	return file(id, loc, "csv")
}

// FileGo defines Writer implementation
// which creates go file
func FileGo(id, loc string) (io.WriteCloser, error) {
	return file(id, loc, "go")
}

// FileGopium defines Writer implementation
// which creates gopium file
func FileGopium(id, loc string) (io.WriteCloser, error) {
	return file(id, loc, "gopium")
}
