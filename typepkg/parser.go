package typepkg

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"

	"1pkg/gopium"

	"golang.org/x/tools/go/packages"
)

// ParserXToolPackagesAst defines packages Parser default implementation
// that uses "golang.org/x/tools/go/packages" packages.Load with cfg to collect package types
// and uses "go/parser" parser.ParseDir to collect ast package
type ParserXToolPackagesAst struct {
	Pattern    string
	AbsDir     string
	ModeTypes  packages.LoadMode
	ModeAst    parser.Mode
	BuildEnv   []string
	BuildFlags []string
}

// ParseTypes ParserXToolPackagesAst implementation
func (p ParserXToolPackagesAst) ParseTypes(ctx context.Context) (*types.Package, gopium.Locator, error) {
	// create packages.Config obj
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Fset:       fset,
		Context:    ctx,
		Dir:        p.AbsDir,
		Mode:       p.ModeTypes,
		Env:        p.BuildEnv,
		BuildFlags: p.BuildFlags,
		Tests:      true,
	}
	// use packages.Load
	pkgs, err := packages.Load(cfg, "")
	// on any error just propagate it
	if err != nil {
		return nil, nil, err
	}
	// check parse results
	// it should be equal to
	// package pattern or last
	// component of the path
	bpkg := path.Base(p.AbsDir)
	if len(pkgs) != 1 ||
		(pkgs[0].String() != p.Pattern && pkgs[0].String() != bpkg) {
		return nil, nil, fmt.Errorf("package %q wasn't found", p.Pattern)
	}
	return pkgs[0].Types, NewLocator(fset), nil
}

// ParseAst ParserXToolPackagesAst implementation
func (p ParserXToolPackagesAst) ParseAst(ctx context.Context) (*ast.Package, gopium.Locator, error) {
	// use parser.ParseDir
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(
		fset,
		p.AbsDir,
		nil,
		p.ModeAst,
	)
	// on any error just propagate it
	if err != nil {
		return nil, nil, err
	}
	// check parse results
	// it should be equal to
	// package pattern or last
	// component of the path
	if pkg, ok := pkgs[p.Pattern]; ok {
		return pkg, NewLocator(fset), nil
	}
	bpkg := path.Base(p.AbsDir)
	if pkg, ok := pkgs[bpkg]; ok {
		return pkg, NewLocator(fset), nil
	}
	return nil, nil, fmt.Errorf("package %q wasn't found", p.Pattern)
}
