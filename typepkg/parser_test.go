package typepkg

import (
	"context"
	"go/parser"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestParserXToolPackagesAstTypesMixed(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("./..")
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p   ParserXToolPackagesAst
		ctx context.Context
		pkg bool
		loc bool
		err bool
	}{
		"non existed folder should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern:   "test",
				Path:      "./test",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			err: true,
		},
		"incorrect pattern with relative path should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern:   "test",
				Path:      "./..",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			err: true,
		},
		"correct pattern with root path should return relevant parser package": {
			p: ParserXToolPackagesAst{
				Pattern:   "1pkg/gopium",
				Root:      pdir,
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: true,
			loc: true,
		},
		"empty types mode should return empty parser package": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium",
				Path:    "./..",
			},
			ctx: context.Background(),
			loc: true,
		},
		"correct pattern and path and mode should return relevant parser package": {
			p: ParserXToolPackagesAst{
				Pattern:   "1pkg/gopium",
				Path:      "./..",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: context.Background(),
			pkg: true,
			loc: true,
		},
		"correct pattern and path and mode should return parser error on on canceled context": {
			p: ParserXToolPackagesAst{
				Pattern:   "1pkg/gopium",
				Path:      "./..",
				ModeTypes: packages.LoadAllSyntax,
			},
			ctx: cctx,
			err: true,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			pkg, loc, err := tcase.p.ParseTypes(tcase.ctx)
			// check
			if tcase.pkg && pkg == nil {
				t.Errorf("actual %v doesn't equal to expected not %v", pkg, nil)
			}
			if !tcase.pkg && pkg != nil {
				t.Errorf("actual %v doesn't equal to expected %v", pkg, nil)
			}
			if tcase.loc && loc == nil {
				t.Errorf("actual %v doesn't equal to expected not %v", loc, nil)
			}
			if !tcase.loc && loc != nil {
				t.Errorf("actual %v doesn't equal to expected %v", loc, nil)
			}
			if tcase.err && err == nil {
				t.Errorf("actual %v doesn't equal to expected not %v", err, nil)
			}
			if !tcase.err && err != nil {
				t.Errorf("actual %v doesn't equal to expected %v", err, nil)
			}
		})
	}
}

func TestParserXToolPackagesAstAstMixed(t *testing.T) {
	// prepare
	pdir, err := filepath.Abs("./..")
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p   ParserXToolPackagesAst
		ctx context.Context
		ast bool
		loc bool
		err bool
	}{
		"non existed folder should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "test",
				Path:    "./test",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			err: true,
		},
		"incorrect pattern with relative path should return parser error": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium",
				Path:    "./..",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			err: true,
		},
		"correct pattern with root path should return relevant parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Root:    pdir,
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			ast: true,
			loc: true,
		},
		"incorrect pattern with full path should return relevant parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "1pkg/gopium",
				Path:    pdir,
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			ast: true,
			loc: true,
		},
		"correct pattern and path and empty ast mode should return relevant parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    "./..",
			},
			ctx: context.Background(),
			ast: true,
			loc: true,
		},
		"correct pattern and path and mode should return relevant parser ast": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    "./..",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: context.Background(),
			ast: true,
			loc: true,
		},
		"correct pattern and path and mode should return parser error on on canceled context": {
			p: ParserXToolPackagesAst{
				Pattern: "gopium",
				Path:    "./..",
				ModeAst: parser.ParseComments | parser.AllErrors,
			},
			ctx: cctx,
			err: true,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			ast, loc, err := tcase.p.ParseAst(tcase.ctx)
			// check
			if tcase.ast && ast == nil {
				t.Errorf("actual %v doesn't equal to expected not %v", ast, nil)
			}
			if !tcase.ast && ast != nil {
				t.Errorf("actual %v doesn't equal to expected %v", ast, nil)
			}
			if tcase.loc && loc == nil {
				t.Errorf("actual %v doesn't equal to expected not %v", loc, nil)
			}
			if !tcase.loc && loc != nil {
				t.Errorf("actual %v doesn't equal to expected %v", loc, nil)
			}
			if tcase.err && err == nil {
				t.Errorf("actual %v doesn't equal to expected not %v", err, nil)
			}
			if !tcase.err && err != nil {
				t.Errorf("actual %v doesn't equal to expected %v", err, nil)
			}
		})
	}
}
