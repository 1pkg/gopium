package typepkg

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"go/types"
	"path"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests"

	"golang.org/x/tools/go/packages"
)

func TestParserXToolPackagesAstTypes(t *testing.T) {
	// prepare
	var wg sync.WaitGroup
	pdir, err := filepath.Abs(filepath.Join("..", "gopium"))
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	sdir, err := filepath.Abs(".")
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p   ParserXToolPackagesAst
		ctx context.Context
		src []byte
		pkg *types.Package
		loc gopium.Locator
		err error
	}{
		"invalid folder should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "test",
				Path:    "test",
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			err: tests.OnOS(
				"windows",
				fmt.Errorf("%s", "couldn't run 'go': chdir test: The system cannot find the file specified."),
				fmt.Errorf("%s", "couldn't run 'go': chdir test: no such file or directory"),
			).(error),
		},
		"invalid pattern with abs path should return expected parser package": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    pdir,
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"empty types mode should return expected empty parser package": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    path.Join("..", "gopium"),
			},
			ctx: context.Background(),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser package": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    path.Join("..", "gopium"),
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser package on abs path": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    pdir,
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser package with tests": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/typepkg",
				Path:    sdir,
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser package skip src": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    path.Join("..", "gopium"),
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			src: []byte(`
//go:build tests_data

package single

type Single struct {
	A	string
	B	string
	C	string
}
`),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return parser error on canceled context": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    path.Join("..", "gopium"),
				//nolint
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: cctx,
			err: context.Canceled,
		},
	}
	for name, tcase := range table {
		// run all parser tests
		// in separate goroutine
		wg.Add(1)
		name := name
		tcase := tcase
		go func(t *testing.T) {
			defer wg.Done()
			t.Run(name, func(t *testing.T) {
				// exec
				pkg, loc, err := tcase.p.ParseTypes(tcase.ctx, tcase.src...)
				// check
				// in case pkg or loc non nil
				// just copy them from result
				if tcase.pkg != nil {
					tcase.pkg = pkg
				}
				if tcase.loc != nil {
					tcase.loc = loc
				}
				if !reflect.DeepEqual(pkg, tcase.pkg) {
					t.Errorf("actual %v doesn't equal to expected %v", pkg, tcase.pkg)
				}
				if !reflect.DeepEqual(loc, tcase.loc) {
					t.Errorf("actual %v doesn't equal to expected %v", loc, tcase.loc)
				}
				if !reflect.DeepEqual(err, tcase.err) {
					// skip the case when error messages are equal
					if (err != nil && tcase.err != nil) &&
						(reflect.DeepEqual(err.Error(), tcase.err.Error())) {
						return
					}
					t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
				}
			})
		}(t)
	}
	// wait util tests finish
	wg.Wait()
}

func TestParserXToolPackagesAstAst(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs(filepath.Join("..", "gopium"))
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p   ParserXToolPackagesAst
		ctx context.Context
		src []byte
		pkg *ast.Package
		loc gopium.Locator
		err error
	}{
		"invalid folder should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "test",
				Path:    "test",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			err: tests.OnOS(
				"windows",
				fmt.Errorf("%s", "open test: The system cannot find the file specified."),
				fmt.Errorf("%s", "open test: no such file or directory"),
			).(error),
		},
		"invalid pattern with relative path should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "github.com/1pkg/gopium/gopium",
				Path:    ".",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			err: errors.New(`ast package "github.com/1pkg/gopium/gopium" wasn't found at "."`),
		},
		"valid pattern with root path should return expected parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Root:    pdir,
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			pkg: &ast.Package{},
			loc: NewLocator(nil),
		},
		"invalid pattern with full path should return expected parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium/1gopium",
				Path:    pdir,
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			pkg: &ast.Package{},
			loc: NewLocator(nil),
		},
		"valid pattern and path and empty ast mode should return expected parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    filepath.Join("..", "gopium"),
			},
			ctx: context.Background(),
			pkg: &ast.Package{},
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    filepath.Join("..", "gopium"),
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			pkg: &ast.Package{},
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser ast with src": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    filepath.Join("..", "gopium"),
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			src: []byte(`
//go:build tests_data

package single

type Single struct {
	A	string
	B	string
	C	string
}
`),
			pkg: &ast.Package{
				Name: "pkg",
				Files: map[string]*ast.File{
					"file": {},
				},
			},
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return parser error with invalid src": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    filepath.Join("..", "gopium"),
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			src: []byte(`
random sets of string
invalid gocode
`),
			err: scanner.ErrorList{
				&scanner.Error{
					Pos: token.Position{Offset: 1, Line: 2, Column: 1},
					Msg: ("expected 'package', found random"),
				},
				&scanner.Error{
					Pos: token.Position{Offset: 13, Line: 2, Column: 13},
					Msg: ("expected ';', found of"),
				},
			},
		},
		"invalid pattern with relative path should return expected parser ast with src": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium/1gopium",
				Path:    filepath.Join("..", "gopium"),
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			src: []byte(`
//go:build tests_data

package single

type Single struct {
	A	string
	B	string
	C	string
}
`),
			pkg: &ast.Package{
				Name: "pkg",
				Files: map[string]*ast.File{
					"file": {},
				},
			},
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return parser error on canceled context": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    filepath.Join("..", "gopium"),
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: cctx,
			err: context.Canceled,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			pkg, loc, err := tcase.p.ParseAst(tcase.ctx, tcase.src...)
			// check
			// in case pkg or loc non nil
			// just copy them from result
			if tcase.pkg != nil {
				tcase.pkg = pkg
			}
			if tcase.loc != nil {
				tcase.loc = loc
			}
			if !reflect.DeepEqual(pkg, tcase.pkg) {
				t.Errorf("actual %v doesn't equal to expected %v", pkg, tcase.pkg)
			}
			if !reflect.DeepEqual(loc, tcase.loc) {
				t.Errorf("actual %v doesn't equal to expected %v", loc, tcase.loc)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				// skip the case when error messages are equal
				if (err != nil && tcase.err != nil) &&
					(reflect.DeepEqual(err.Error(), tcase.err.Error())) {
					return
				}
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
