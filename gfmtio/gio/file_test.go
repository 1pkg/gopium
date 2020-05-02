package gio

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("./..")
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		name string
		path string
		ext  string
		full string
		err  error
	}{
		"should create file with simple valid params": {
			name: "test",
			path: pdir,
			ext:  "test",
			full: path.Join(filepath.Dir(pdir), "test.test"),
		},
		"should create same file with long name param": {
			name: "test/test/test.test",
			path: pdir,
			ext:  "test",
			full: path.Join(filepath.Dir(pdir), "test.test"),
		},
		"should create different file with long name param and long extension": {
			name: "test/test/test.test",
			path: pdir,
			ext:  "test.test",
			full: path.Join(filepath.Dir(pdir), "test.test.test"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			defer os.Remove(tcase.full)
			wc, err := file(tcase.name, tcase.path, tcase.ext)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			if err == nil && reflect.DeepEqual(wc, nil) {
				t.Errorf("actual %v doesn't equal to expected not %v", wc, nil)
			}
			if _, err := os.Stat(tcase.full); err != nil {
				t.Errorf("actual %v doesn't equal to expected %v", err, nil)
			}
		})
	}
}
