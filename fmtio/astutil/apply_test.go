package astutil

import (
	"bytes"
	"context"
	"errors"
	"go/parser"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
	"1pkg/gopium/gopium"
	"1pkg/gopium/tests"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
	"1pkg/gopium/typepkg"
)

func TestApply(t *testing.T) {
	// prepare
	lh := collections.NewHierarchic(tests.Gopium)
	lh.Push(
		"tests_data_note_file-1.go:6",
		filepath.Join(tests.Gopium, "tests", "data", "note", "file-1.go"),
		gopium.Struct{
			Name:    "Note",
			Doc:     []string{"// test-doc", "// test-doc-doc"},
			Comment: []string{"// test-com"},
			Fields: []gopium.Field{
				{
					Name: "C",
					Type: "string",
				},
				{
					Name: "_",
					Type: "[]byte",
					Size: 8,
					Doc:  []string{"// test-pad", "// test-pad-pad"},
				},
				{
					Name: "A",
					Type: "string",
				},
			},
		},
	)
	lh.Push(
		"tests_data_note_file-2.go:6",
		filepath.Join(tests.Gopium, "tests", "data", "note", "file-2.go"),
		gopium.Struct{
			Name: "DocCom",
			Fields: []gopium.Field{
				{
					Name:    "f",
					Type:    "complex128",
					Comment: []string{"// f com 1", "// f com 2", "// f com 3"},
				},
			},
		},
	)
	ldc := collections.NewHierarchic(tests.Gopium)
	ldc.Push(
		"tests_data_note_file-1.go:6",
		filepath.Join(tests.Gopium, "tests", "data", "note", "file-1.go"),
		gopium.Struct{
			Name: "Note",
			Fields: []gopium.Field{
				{
					Name: "C",
					Type: "string",
				},
				{
					Name: "_",
					Type: "[]byte",
					Size: 8,
				},
				{
					Name: "A",
					Type: "string",
				},
			},
		},
	)
	ldc.Push(
		"tests_data_note_file-2.go:6",
		filepath.Join(tests.Gopium, "tests", "data", "note", "file-2.go"),
	)
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := fmtio.Gofmt{}
	sp := Package{}
	table := map[string]struct {
		p   gopium.Parser
		a   gopium.Apply
		ctx context.Context
		h   collections.Hierarchic
		r   map[string][]byte
		err error
	}{
		"empty pkg should apply nothing": {
			p:   data.NewParser("empty"),
			a:   UFFN,
			ctx: context.Background(),
			r:   map[string][]byte{},
		},
		"note struct pkg should apply expected structs": {
			p:   data.NewParser("note"),
			a:   UFFN,
			ctx: context.Background(),
			h:   lh,
			r: map[string][]byte{
				"tests_data_note_file-1.go": []byte(`
//+build tests_data

package note

// Note doc
// test-doc test-doc-doc
type Note struct {
	C string
	// test-pad test-pad-pad
	_ [8]byte
	A string
} // test-com
// some comment

// last comment
`),
				"tests_data_note_file-2.go": []byte(`
//+build tests_data

package note

type DocCom struct {
	f complex128 // f com 1 f com 2 f com 3
} // doc com
`),
			},
		},
		"note struct pkg should skip expected structs": {
			p:   data.NewParser("note"),
			a:   UFFN,
			ctx: context.Background(),
			h:   ldc,
			r: map[string][]byte{
				"tests_data_note_file-1.go": []byte(`
//+build tests_data

package note

// Note doc
type Note struct {
	C string
	_ [8]byte
	A string
} // some comment

// last comment
`),
				"tests_data_note_file-2.go": []byte(`
//+build tests_data

package note

/* 1pkg - 🌺 gopium @1pkg */
type DocCom struct {
	f complex128
	// 🌺 gopium @1pkg
} // doc com
`),
			},
		},
		"note struct pkg should apply nothing on canceled context": {
			p:   data.NewParser("note"),
			a:   UFFN,
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context fast": {
			p:   data.NewParser("note"),
			a:   ufmt(walk, mocks.Ast{}.Ast),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context filter": {
			p:   data.NewParser("note"),
			a:   filter(walk),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context filter after walk": {
			p:   data.NewParser("note"),
			a:   filter(mocks.Walk{}.Walk),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on walk error filter": {
			p:   data.NewParser("note"),
			a:   filter(mocks.Walk{Err: errors.New("walk-test")}.Walk),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("walk-test"),
		},
		"note struct pkg should apply nothing on canceled context note": {
			p: data.NewParser("note"),
			a: note(
				walk,
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				fmtio.Gofmt{},
			),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context after walk note": {
			p: data.NewParser("note"),
			a: note(
				mocks.Walk{Err: errors.New("walk-test")}.Walk,
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				fmtio.Gofmt{},
			),
			h:   lh,
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("walk-test"),
		},
		"note struct pkg should apply nothing on parser error": {
			p: data.NewParser("note"),
			a: note(
				walk,
				mocks.Parser{Asterr: errors.New("test-1")},
				fmtio.Gofmt{},
			),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("test-1"),
		},
		"note struct pkg should apply nothing on printer error": {
			p: data.NewParser("note"),
			a: note(
				walk,
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				mocks.Printer{Err: errors.New("test-2")},
			),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("test-2"),
		},
		"note struct pkg should apply nothing on apply error": {
			p: data.NewParser("note"),
			a: combine(
				mocks.Apply{Err: errors.New("test-3")}.Apply,
				filter(walk),
			),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("test-3"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			w := &mocks.Writer{}
			pkg, loc, err := tcase.p.ParseAst(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
			}
			// exec
			pkg, err = tcase.a(tcase.ctx, pkg, loc, tcase.h)
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// prepare
			if pkg != nil {
				err = sp.Persist(context.Background(), p, data.Writer{Writer: w}, loc, pkg)
				if !reflect.DeepEqual(err, nil) {
					t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
				}
			}
			// check
			for name, rwc := range w.RWCs {
				// check all struct
				// against bytes map
				if st, ok := tcase.r[name]; ok {
					// read rwc to buffer
					var buf bytes.Buffer
					_, err := buf.ReadFrom(rwc)
					if !reflect.DeepEqual(err, nil) {
						t.Errorf("actual %v doesn't equal to expected %v", err, nil)
					}
					// format actual and expected identically
					actual := strings.Trim(string(buf.Bytes()), "\n")
					expected := strings.Trim(string(st), "\n")
					if !reflect.DeepEqual(actual, expected) {
						t.Errorf("name %v actual %v doesn't equal to expected %v", name, actual, expected)
					}
					delete(tcase.r, name)
				} else {
					t.Errorf("actual %v doesn't equal to expected %v", name, "")
				}
			}
			// check that map has been drained
			if !reflect.DeepEqual(tcase.r, map[string][]byte{}) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.r, map[string][]byte{})
			}
		})
	}
}
