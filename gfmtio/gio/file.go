package gio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// file defines Writer helper
// which creates underlying file
// by name, path and ext
func file(name, path, ext string) (io.WriteCloser, error) {
	bname := filepath.Base(name)
	bname = strings.Split(bname, ".")[0]
	path = filepath.Dir(path)
	return os.Create(fmt.Sprintf("%s/%s.%s", path, bname, ext))
}
