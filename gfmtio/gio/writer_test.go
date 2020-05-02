package gio

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriter(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("./..")
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		w    Writer
		id   string
		loc  string
		full string
		err  error
	}{
		"stdout should return stdout back": {
			w:    Stdout,
			id:   "test",
			loc:  pdir,
			full: "",
		},
		"filejson should create valid json file": {
			w:    FileJson,
			id:   "test",
			loc:  pdir,
			full: path.Join(filepath.Dir(pdir), "test.json"),
		},
		"filexml should create valid xml file": {
			w:    FileXml,
			id:   "test",
			loc:  pdir,
			full: path.Join(filepath.Dir(pdir), "test.xml"),
		},
		"filecs should create valid csv file": {
			w:    FileCsv,
			id:   "test",
			loc:  pdir,
			full: path.Join(filepath.Dir(pdir), "test.csv"),
		},
		"filego should create valid go file": {
			w:    FileGo,
			id:   "test",
			loc:  pdir,
			full: path.Join(filepath.Dir(pdir), "test.go"),
		},
		"filegopium should create valid gopium file": {
			w:    FileGopium,
			id:   "test",
			loc:  pdir,
			full: path.Join(filepath.Dir(pdir), "test.gopium"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			if tcase.full != "" {
				defer os.Remove(tcase.full)
			}
			wc, err := tcase.w(tcase.id, tcase.loc)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			if err == nil && reflect.DeepEqual(wc, nil) {
				t.Errorf("actual %v doesn't equal to expected not %v", wc, nil)
			}
			if _, err := os.Stat(tcase.full); tcase.full != "" && err != nil {
				t.Errorf("actual %v doesn't equal to expected %v", err, nil)
			}
		})
	}
}
