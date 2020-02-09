package pkgs

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// Parser defines abstraction for packages parsing processor
type Parser interface {
	Parse(context.Context) (*types.Package, *ast.Package, error)
}

// ParserXTool defines packages Parser default "golang.org/x/tools/go/packages" implementation
// that uses packages.Load with cfg to collect package types
type ParserXToolPackages struct {
	Pattern    string
	AbsDir     string
	LoadMode   packages.LoadMode
	BuildEnv   []string
	BuildFlags []string
}

// Parse packages Parser default "golang.org/x/tools/go/packages" implementation
func (p ParserXToolPackages) Parse(ctx context.Context) (*types.Package, *ast.Package, error) {
	// create packages.Config obj
	cfg := &packages.Config{
		Fset:       token.NewFileSet(),
		Context:    ctx,
		Dir:        p.AbsDir,
		Mode:       p.LoadMode,
		Env:        p.BuildEnv,
		BuildFlags: p.BuildFlags,
		Tests:      true,
	}
	// use load packages
	pkgs, err := packages.Load(cfg, p.Pattern)
	if err != nil {
		return nil, nil, err
	}
	// check parse results
	if len(pkgs) != 1 || pkgs[0].String() != p.Pattern {
		return nil, nil, fmt.Errorf("packages %q wasn't found", p.Pattern)
	}
	return pkgs[0].Types, nil, nil
}

// ParserMock defines packages Parser mock implementation
type ParserMock struct {
	tpkg *types.Package
	apkg *ast.Package
}

// Parse ParserMock implementation
func (p ParserMock) Parse(context.Context) (*types.Package, *ast.Package, error) {
	return p.tpkg, p.apkg, nil
}

// ParserNil defines packages Parser nil implementation
type ParserNil struct{}

// Parse ParserNil implementation
func (ParserNil) Parse(context.Context) (*types.Package, *ast.Package, error) {
	return nil, nil, nil
}

// ParserError defines packages Parser error implementation
type ParserError struct {
	err error
}

// Parse ParserError implementation
func (p ParserError) Parse(context.Context) (*types.Package, *ast.Package, error) {
	return nil, nil, p.err
}
