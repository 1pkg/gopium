package fmtio

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"1pkg/gopium/gopium"
)

func TestWriter(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("..")
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pfile, err := filepath.Abs(filepath.Join("..", "opium.go"))
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		w    gopium.Writer
		loc  string
		path string
		err  error
		werr error
		cerr error
	}{
		"stdout should return expected stdout writer": {
			w:    Stdout{},
			loc:  pdir,
			path: "",
		},
		"file json should return expected json writer": {
			w:    File{Name: "test", Ext: "json"},
			loc:  pfile,
			path: filepath.Join(pdir, "test.json"),
		},
		"file xml should return expected xml writer": {
			w:    File{Name: "test", Ext: "xml"},
			loc:  pfile,
			path: filepath.Join(pdir, "test.xml"),
		},
		"file csv should return expected csv writer": {
			w:    File{Name: "test", Ext: "csv"},
			loc:  pfile,
			path: filepath.Join(pdir, "test.csv"),
		},
		"files json should return expected json writer": {
			w:    Files{Ext: "json"},
			loc:  pfile,
			path: filepath.Join(pdir, "opium.json"),
		},
		"files xml should return expected xml writer": {
			w:    Files{Ext: "xml"},
			loc:  pfile,
			path: filepath.Join(pdir, "opium.xml"),
		},
		"files csv should return expected csv writer": {
			w:    Files{Ext: "csv"},
			loc:  pfile,
			path: filepath.Join(pdir, "opium.csv"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			wc, err := tcase.w.Generate(tcase.loc)
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
			if tcase.path != "" {
				defer os.Remove(tcase.path)
				if _, err := os.Stat(tcase.path); !reflect.DeepEqual(err, nil) {
					t.Errorf("actual %v doesn't equal to expected %v", err, nil)
				}
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

func TestCategoryWriter(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("..")
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pfile, err := filepath.Abs(filepath.Join("..", "opium.go"))
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		w    gopium.CategoryWriter
		cat  string
		loc  string
		path string
		err  error
	}{
		"file json with replace cat writer should return expected json writer": {
			w:    Origin{Writter: File{Name: "test", Ext: "json"}},
			cat:  pdir,
			loc:  pfile,
			path: filepath.Join(pdir, "test.json"),
		},
		"file json with copy cat writer should return expected json writer": {
			w:    &Suffix{Writter: File{Name: "test", Ext: "json"}, Suffix: "test"},
			cat:  pdir,
			loc:  pfile,
			path: filepath.Join(fmt.Sprintf("%s_%s", pdir, "test"), "test.json"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			err := tcase.w.Category(tcase.cat)
			wc, gerr := tcase.w.Generate(tcase.loc)
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
			if tcase.path != "" {
				defer func() {
					os.Remove(tcase.path)
					os.Remove(filepath.Dir(tcase.path))
				}()
				if _, err := os.Stat(tcase.path); !reflect.DeepEqual(err, nil) {
					t.Errorf("actual %v doesn't equal to expected %v", err, nil)
				}
			}
			if !reflect.DeepEqual(gerr, nil) {
				t.Errorf("actual %v doesn't equal to expected %v", gerr, nil)
			}
			if !reflect.DeepEqual(werr, nil) {
				t.Errorf("actual %v doesn't equal to expected %v", werr, nil)
			}
			if !reflect.DeepEqual(n, 0) {
				t.Errorf("actual %v doesn't equal to expected %v", n, 0)
			}
			if !reflect.DeepEqual(cerr, nil) {
				t.Errorf("actual %v doesn't equal to expected %v", werr, nil)
			}
		})
	}
}
