package typepkg

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"1pkg/gopium"

	"golang.org/x/tools/go/packages"
)

// ParserXToolPackagesAst defines
// gopium parser default implementation
// that uses "golang.org/x/tools/go/packages"
// to collect package types
// and "go/parser" to collect ast package
type ParserXToolPackagesAst struct {
	Pattern    string
	Path       string
	Root       string
	ModeTypes  packages.LoadMode
	ModeAst    parser.Mode
	BuildEnv   []string
	BuildFlags []string
}

// ParseTypes ParserXToolPackagesAst implementation
func (p ParserXToolPackagesAst) ParseTypes(ctx context.Context) (*types.Package, gopium.Locator, error) {
	// create packages.Config obj
	fset := token.NewFileSet()
	dir := filepath.Join(p.Root, p.Path)
	cfg := &packages.Config{
		Fset:       fset,
		Context:    ctx,
		Dir:        dir,
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
	// package pattern or
	// all except first components of path
	pkg := p.Path
	if list := strings.Split(p.Path, string(filepath.Separator)); len(list) > 0 {
		pkg = filepath.Join(list[1:]...)
	}
	if len(pkgs) != 1 ||
		(pkgs[0].String() != p.Pattern && pkgs[0].String() != pkg) {
		return nil, nil, fmt.Errorf("package %q wasn't found at %q", p.Pattern, dir)
	}
	return pkgs[0].Types, NewLocator(fset), nil
}

// ParseAst ParserXToolPackagesAst implementation
func (p ParserXToolPackagesAst) ParseAst(ctx context.Context) (*ast.Package, gopium.Locator, error) {
	// manage context actions
	// in case of cancelation
	// stop parse and return error back
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}
	// use parser.ParseDir
	fset := token.NewFileSet()
	dir := filepath.Join(p.Root, p.Path)
	pkgs, err := parser.ParseDir(
		fset,
		dir,
		nil,
		p.ModeAst,
	)
	// on any error just propagate it
	if err != nil {
		return nil, nil, err
	}
	// check parse results
	// it should be equal to
	// package pattern or
	// last component of path
	if pkg, ok := pkgs[p.Pattern]; len(pkgs) == 1 && ok {
		return pkg, NewLocator(fset), nil
	}
	pkg := filepath.Base(p.Path)
	if pkg, ok := pkgs[pkg]; len(pkgs) == 1 && ok {
		return pkg, NewLocator(fset), nil
	}
	return nil, nil, fmt.Errorf("package %q wasn't found at %q", p.Pattern, dir)
}
