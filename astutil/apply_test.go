package astutil

import (
	"bytes"
	"context"
	"errors"
	"go/build"
	"go/parser"
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

func TestApply(t *testing.T) {
	// prepare
	lh := collections.NewHierarchic(build.Default.GOPATH)
	lh.Push(
		"6-39ba0c31867d8eaabd59a515e15955bbe83b4aa800278c7ef0c75e5ca9bcf56c",
		"/src/1pkg/gopium/tests/data/note/file-1.go",
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
		"6-90fba0480e71f274086a3057fe48a45c98599132b3e64b02d2b7540bb385e217",
		"/src/1pkg/gopium/tests/data/note/file-2.go",
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
	ldc := collections.NewHierarchic(build.Default.GOPATH)
	ldc.Push(
		"6-39ba0c31867d8eaabd59a515e15955bbe83b4aa800278c7ef0c75e5ca9bcf56c",
		"/src/1pkg/gopium/tests/data/note/file-1.go",
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
		"6-90fba0480e71f274086a3057fe48a45c98599132b3e64b02d2b7540bb385e217",
		"/src/1pkg/gopium/tests/data/note/file-2.go",
	)
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := fmtio.Goprint(0, 4, false)
	table := map[string]struct {
		p   gopium.Parser
		a   Apply
		ctx context.Context
		h   collections.Hierarchic
		r   map[string][]byte
		err error
	}{
		"empty pkg should apply nothing": {
			p:   data.NewParser("empty"),
			a:   UFFN,
			ctx: context.Background(),
			r: map[string][]byte{
				"/src/1pkg/gopium/tests/data/empty/file.go": []byte(`
//+build tests_data

package empty
`),
			},
		},
		"note struct pkg should apply expected structs": {
			p:   data.NewParser("note"),
			a:   UFFN,
			ctx: context.Background(),
			h:   lh,
			r: map[string][]byte{
				"/src/1pkg/gopium/tests/data/note/file-1.go": []byte(`
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
				"/src/1pkg/gopium/tests/data/note/file-2.go": []byte(`
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
				"/src/1pkg/gopium/tests/data/note/file-1.go": []byte(`
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
				"/src/1pkg/gopium/tests/data/note/file-2.go": []byte(`
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
			r:   make(map[string][]byte),
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context fast": {
			p:   data.NewParser("note"),
			a:   ufmt(walk, mocks.Ast{}.Ast),
			ctx: cctx,
			r:   make(map[string][]byte),
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context filter": {
			p:   data.NewParser("note"),
			a:   filter(walk),
			ctx: cctx,
			r:   make(map[string][]byte),
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context filter after walk": {
			p:   data.NewParser("note"),
			a:   filter(mocks.Walk{}.Walk),
			ctx: cctx,
			r:   make(map[string][]byte),
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on walk error filter": {
			p:   data.NewParser("note"),
			a:   filter(mocks.Walk{Err: errors.New("walk-test")}.Walk),
			ctx: context.Background(),
			r:   make(map[string][]byte),
			err: errors.New("walk-test"),
		},
		"note struct pkg should apply nothing on canceled context note": {
			p: data.NewParser("note"),
			a: note(
				walk,
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				fmtio.Goprint(0, 4, false),
			),
			ctx: cctx,
			r:   make(map[string][]byte),
			err: context.Canceled,
		},
		"note struct pkg should apply nothing on canceled context after walk note": {
			p: data.NewParser("note"),
			a: note(
				mocks.Walk{Err: errors.New("walk-test")}.Walk,
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				},
				fmtio.Goprint(0, 4, false),
			),
			h:   lh,
			ctx: context.Background(),
			r:   make(map[string][]byte),
			err: errors.New("walk-test"),
		},
		"note struct pkg should apply nothing on parser error": {
			p: data.NewParser("note"),
			a: note(
				walk,
				mocks.Parser{Asterr: errors.New("test-1")},
				fmtio.Goprint(0, 4, false),
			),
			ctx: context.Background(),
			r:   make(map[string][]byte),
			err: errors.New("test-1"),
		},
		"note struct pkg should apply nothing on printer error": {
			p: data.NewParser("note"),
			a: note(
				walk,
				&typepkg.ParserXToolPackagesAst{
					ModeAst: parser.ParseComments | parser.AllErrors,
				}, mocks.Printer{Err: errors.New("test-2")}.Printer,
			),
			ctx: context.Background(),
			r:   make(map[string][]byte),
			err: errors.New("test-2"),
		},
		"note struct pkg should apply nothing on apply error": {
			p: data.NewParser("note"),
			a: combine(
				mocks.Apply{Err: errors.New("test-3")}.Apply,
				filter(walk),
			),
			ctx: context.Background(),
			r:   make(map[string][]byte),
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
				err = p.Save(w.Writer)(context.Background(), apkg, loc)
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
					buf.ReadFrom(rwc)
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
			if !reflect.DeepEqual(tcase.r, make(map[string][]byte)) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.r, make(map[string][]byte))
			}
		})
	}
}
