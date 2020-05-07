package fmtio

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
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		w    Writer
		id   string
		loc  string
		path string
		err  error
		werr error
		cerr error
	}{
		"stdout should return expected stdout writer": {
			w:    Stdout,
			id:   "test",
			loc:  pdir,
			path: "",
		},
		"filejson should return expected json writer": {
			w:    File("json"),
			id:   "test",
			loc:  pdir,
			path: path.Join(filepath.Dir(pdir), "test.json"),
		},
		"filexml should return expected xml writer": {
			w:    File("xml"),
			id:   "test",
			loc:  pdir,
			path: path.Join(filepath.Dir(pdir), "test.xml"),
		},
		"filecs should return expected csv writer": {
			w:    File("csv"),
			id:   "test",
			loc:  pdir,
			path: path.Join(filepath.Dir(pdir), "test.csv"),
		},
		"filego should return expected go writer": {
			w:    File("go"),
			id:   "test",
			loc:  pdir,
			path: path.Join(filepath.Dir(pdir), "test.go"),
		},
		"filegopium should return expected gopium writer": {
			w:    File("gopium"),
			id:   "test",
			loc:  pdir,
			path: path.Join(filepath.Dir(pdir), "test.gopium"),
		},
		"long id param should return expected writer": {
			w:    File("test"),
			id:   "test/test/test.test",
			loc:  pdir,
			path: path.Join(filepath.Dir(pdir), "test.test"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			if tcase.path != "" {
				defer os.Remove(tcase.path)
			}
			// exec
			wc, err := tcase.w(tcase.id, tcase.loc)
			n, werr := wc.Write([]byte(``))
			cerr := wc.Close()
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			if err == nil && reflect.DeepEqual(wc, nil) {
				t.Errorf("actual %v doesn't equal to expected not %v", wc, nil)
			}
			// check that such file exists
			if _, err := os.Stat(tcase.path); tcase.path != "" && !reflect.DeepEqual(err, nil) {
				t.Errorf("actual %v doesn't equal to expected %v", err, nil)
			}
			if !reflect.DeepEqual(werr, tcase.werr) {
				t.Errorf("actual %v doesn't equal to expected %v", werr, tcase.werr)
			}
			if !reflect.DeepEqual(n, 0) {
				t.Errorf("actual %v doesn't equal to expected %v", n, 0)
			}
			if !reflect.DeepEqual(cerr, tcase.cerr) {
				t.Errorf("actual %v doesn't equal to expected %v", werr, tcase.werr)
			}
		})
	}
}
