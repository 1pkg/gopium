package typepkg

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"
	"testing"

	"1pkg/gopium"

	"golang.org/x/tools/go/packages"
)

func TestParserXToolPackagesAstTypesMixed(t *testing.T) {
	// prepare
	var wg sync.WaitGroup
	pdir, err := filepath.Abs("./..")
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p   ParserXToolPackagesAst
		ctx context.Context
		pkg *types.Package
		loc gopium.Locator
		err error
	}{
		"invalid folder should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern:   "test",
				Path:      "./test",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			err: errors.New("couldn't exec 'go [-e -json -compiled=true -test=true -export=false -deps=true -find=false -- ]': chdir test: no such file or directory *os.PathError"),
		},
		"invalid pattern with relative path should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern:   "test",
				Path:      "./..",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			err: errors.New(`package "test" wasn't found at ".."`),
		},
		"invalid pattern with root path should return expected parser package": {
			p: ParserXToolPackagesAst{
				Pattern:   "1pkg/gopium",
				Root:      pdir,
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"empty types mode should return expected empty parser package": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium",
				Path:    "./..",
			},
			ctx: context.Background(),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser package": {
			p: ParserXToolPackagesAst{
				Pattern:   "1pkg/gopium",
				Path:      "./..",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: types.NewPackage("test", "test"),
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return parser error on canceled context": {
			p: ParserXToolPackagesAst{
				Pattern:   "1pkg/gopium",
				Path:      "./..",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: cctx,
			err: cctx.Err(),
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
				pkg, loc, err := tcase.p.ParseTypes(tcase.ctx)
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
					t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
				}
			})
		}(t)
	}
	// wait util tests finish
	wg.Wait()
}

func TestParserXToolPackagesAstAstMixed(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("./..")
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p   ParserXToolPackagesAst
		ctx context.Context
		pkg *ast.Package
		loc gopium.Locator
		err error
	}{
		"invalid folder should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "test",
				Path:    "./test",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			err: &os.PathError{Op: "open", Path: "test", Err: syscall.Errno(2)},
		},
		"invalid pattern with relative path should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium",
				Path:    "./..",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			err: errors.New(`package "1pkg/gopium" wasn't found at ".."`),
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
				Pattern: "1pkg/gopium",
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
				Path:    "./..",
			},
			ctx: context.Background(),
			pkg: &ast.Package{},
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return expected parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    "./..",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			pkg: &ast.Package{},
			loc: NewLocator(nil),
		},
		"valid pattern and path and mode should return parser error on canceled context": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    "./..",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: cctx,
			err: cctx.Err(),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			pkg, loc, err := tcase.p.ParseAst(tcase.ctx)
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
				t.Errorf("actual %+v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
