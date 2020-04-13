package typepkg

import (
	"context"
	"go/parser"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestParserXToolPackagesAstTypesMixed(t *testing.T) {
	// non existed folder
	p := ParserXToolPackagesAst{
		Pattern:   "test",
		Path:      "./test",
		ModeTypes: packages.LoadAllSyntax,
	}
	pkg, loc, err := p.ParseTypes(context.Background())
	if pkg != nil {
		t.Errorf("actual %v doesn't equal to expected %v", pkg, nil)
	}
	if loc != nil {
		t.Errorf("actual %v doesn't equal to expected %v", loc, nil)
	}
	if err == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", err, nil)
	}
	// incorrect patter name
	p = ParserXToolPackagesAst{
		Pattern:   "test",
		Path:      "./..",
		ModeTypes: packages.LoadAllSyntax,
	}
	pkg, loc, err = p.ParseTypes(context.Background())
	if pkg != nil {
		t.Errorf("actual %v doesn't equal to expected %v", pkg, nil)
	}
	if loc != nil {
		t.Errorf("actual %v doesn't equal to expected %v", loc, nil)
	}
	if err == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", err, nil)
	}
	// empty types mode
	p = ParserXToolPackagesAst{
		Pattern: "1pkg/gopium",
		Path:    "./..",
	}
	pkg, loc, err = p.ParseTypes(context.Background())
	if pkg != nil {
		t.Errorf("actual %v doesn't equal to expected %v", pkg, nil)
	}
	if loc == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", loc, nil)
	}
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
	// correct package
	p = ParserXToolPackagesAst{
		Pattern:   "1pkg/gopium",
		Path:      "./..",
		ModeTypes: packages.LoadAllSyntax,
	}
	pkg, loc, err = p.ParseTypes(context.Background())
	if pkg == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", pkg, nil)
	}
	if loc == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", loc, nil)
	}
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
}

func TestParserXToolPackagesAstAstMixed(t *testing.T) {
	// non existed folder
	p := ParserXToolPackagesAst{
		Pattern: "test",
		Path:    "./test",
		ModeAst: parser.ParseComments | parser.AllErrors,
	}
	ast, loc, err := p.ParseAst(context.Background())
	if ast != nil {
		t.Errorf("actual %v doesn't equal to expected %v", ast, nil)
	}
	if loc != nil {
		t.Errorf("actual %v doesn't equal to expected %v", loc, nil)
	}
	if err == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", err, nil)
	}
	// incorrect pattern package with relative path
	p = ParserXToolPackagesAst{
		Pattern: "1pkg/gopium",
		Path:    "./..",
		ModeAst: parser.ParseComments | parser.AllErrors,
	}
	ast, loc, err = p.ParseAst(context.Background())
	if ast != nil {
		t.Errorf("actual %v doesn't equal to expected %v", ast, nil)
	}
	if loc != nil {
		t.Errorf("actual %v doesn't equal to expected %v", loc, nil)
	}
	if err == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", err, nil)
	}
	// incorrect pattern package with full path
	pdir, _ := filepath.Abs("./..")
	p = ParserXToolPackagesAst{
		Pattern: "1pkg/gopium",
		Path:    pdir,
		ModeAst: parser.ParseComments | parser.AllErrors,
	}
	ast, loc, err = p.ParseAst(context.Background())
	if ast == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", ast, nil)
	}
	if loc == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", loc, nil)
	}
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
	// correct pattern package
	p = ParserXToolPackagesAst{
		Pattern: "gopium",
		Path:    "./..",
		ModeAst: parser.ParseComments | parser.AllErrors,
	}
	ast, loc, err = p.ParseAst(context.Background())
	if ast == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", ast, nil)
	}
	if loc == nil {
		t.Errorf("actual %v doesn't equal to expected not %v", loc, nil)
	}
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
}
