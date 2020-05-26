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

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
	"1pkg/gopium/typepkg"
)

func TestXapply(t *testing.T) {
	// prepare
	lh := collections.NewHierarchic(gopium.Root())
	lh.Push(
		"660d36c978f943d2e8325462c049cf1e003521b3ad3fc2f71c646cbf51a3acc1:6",
		filepath.Join(gopium.Root(), "tests", "data", "note", "file-1.go"),
		gopium.Struct{
			Name:    "Note",
			Doc:     []string{"// test-doc"},
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
					Doc:  []string{"// test-pad"},
				},
				{
					Name: "A",
					Type: "string",
				},
			},
		},
	)
	lh.Push(
		"dcd36f56fb9252fc90eac010290a5ae42b67d55ad3c8fbe55a1aa72749633e0e:6",
		filepath.Join(gopium.Root(), "tests", "data", "note", "file-2.go"),
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
	ldc := collections.NewHierarchic(gopium.Root())
	ldc.Push(
		"660d36c978f943d2e8325462c049cf1e003521b3ad3fc2f71c646cbf51a3acc1:6",
		filepath.Join(gopium.Root(), "tests", "data", "note", "file-1.go"),
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
		"dcd36f56fb9252fc90eac010290a5ae42b67d55ad3c8fbe55a1aa72749633e0e:6",
		filepath.Join(gopium.Root(), "tests", "data", "note", "file-2.go"),
	)
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := fmtio.NewGoprinter(0, 4, false)
	sp := Package{}
	table := map[string]struct {
		p   gopium.Parser
		a   gopium.Xapply
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
// test-doc
type Note struct {
	C	string
	// test-pad
	_	[8]byte
	A	string
}	// test-com
// some comment

// last comment
`),
				"tests_data_note_file-2.go": []byte(`
//+build tests_data

package note

/**/
type DocCom struct {
	f complex128	// f com 1
	// f com 2
	// f com 3
}	// doc com
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
	C	string
	_	[8]byte
	A	string
}	// some comment

// last comment
`),
				"tests_data_note_file-2.go": []byte(`
//+build tests_data

package note

/**/
type DocCom struct {
	f complex128
}	// doc com
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
			a:   ufmt(walker{}, mocks.Xast{}.Ast),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context filter": {
			p:   data.NewParser("note"),
			a:   filter(walker{}),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context filter after walk": {
			p:   data.NewParser("note"),
			a:   filter(mocks.Xwalker{}),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on walk error filter": {
			p:   data.NewParser("note"),
			a:   filter(mocks.Xwalker{Err: errors.New("walk-test")}),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("walk-test"),
		},
		"note struct pkg should apply nothing on canceled context note": {
			p: data.NewParser("note"),
			a: note(
				walker{},
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				fmtio.NewGoprinter(0, 4, false),
			),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context after walk note": {
			p: data.NewParser("note"),
			a: note(
				mocks.Xwalker{Err: errors.New("walk-test")},
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				fmtio.NewGoprinter(0, 4, false),
			),
			h:   lh,
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("walk-test"),
		},
		"note struct pkg should apply nothing on parser error": {
			p: data.NewParser("note"),
			a: note(
				walker{},
				mocks.Parser{Asterr: errors.New("test-1")},
				fmtio.NewGoprinter(0, 4, false),
			),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("test-1"),
		},
		"note struct pkg should apply nothing on printer error": {
			p: data.NewParser("note"),
			a: note(
				walker{},
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				}, mocks.Printer{Err: errors.New("test-2")},
			),
			ctx: context.Background(),
			r:   map[string][]byte{},
			err: errors.New("test-2"),
		},
		"note struct pkg should apply nothing on apply error": {
			p: data.NewParser("note"),
			a: combine(
				mocks.Xapply{Err: errors.New("test-3")}.Apply,
				filter(walker{}),
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
			apkg, aerr := tcase.a(tcase.ctx, pkg, loc, tcase.h)
			// prepare
			if apkg != nil {
				err = sp.Persist(context.Background(), p, w, loc, apkg)
				if !reflect.DeepEqual(err, nil) {
					t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
				}
			}
			// check
			if !reflect.DeepEqual(aerr, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", aerr, tcase.err)
			}
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
