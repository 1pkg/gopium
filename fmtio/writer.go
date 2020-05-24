package fmtio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"1pkg/gopium"
)

// Stdout defines writer implementation
// which only returns os stdout all the time
func Stdout(string) (io.WriteCloser, error) {
	return stdout{}, nil
}

// File defines writer helper
// which creates underlying file callback
// with captured name and ext on provided loc
// that implementats writer
func File(name string, ext string) gopium.Writer {
	return func(loc string) (io.WriteCloser, error) {
		path := filepath.Dir(loc)
		return os.Create(fmt.Sprintf("%s/%s.%s", path, name, ext))
	}
}

// Files defines writer helper
// which creates underlying files callback
// with captured ext on provided loc
// that implementats writer
func Files(ext string) gopium.Writer {
	return func(loc string) (io.WriteCloser, error) {
		path := strings.Replace(loc, filepath.Ext(loc), fmt.Sprintf(".%s", ext), 1)
		return os.Create(path)
	}
}

// Replace defines cat writer implementation
// which just returs provided writer back
func Replace(w gopium.Writer, rcat string) (gopium.Writer, error) {
	return w, nil
}

// Gopium defines cat writer implementation
// which replaces root cat for writer with provided suffixed root cat
func Copy(suffix string) gopium.Catwriter {
	return func(w gopium.Writer, rcat string) (gopium.Writer, error) {
		// add suffix to root cat
		sufrcat := fmt.Sprintf("%s_%s", rcat, suffix)
		// create dir for new suffixed cat
		err := os.MkdirAll(sufrcat, os.ModePerm)
		return func(loc string) (io.WriteCloser, error) {
			loc = strings.Replace(loc, rcat, sufrcat, 1)
			return w(loc)
		}, err
	}
}
