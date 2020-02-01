package pkgs

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"regexp"

	"golang.org/x/tools/go/packages"
)

// Parser defines abstraction for package parsing processor
type Parser func(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error)

// ParserXTool defines package parser default go "golang.org/x/tools/go/packages" implementation
// that uses packages load with cfg to collect types, fileset and err
type ParserXTool struct {
	Patterns   []string
	AbsDir     string
	LoadMode   packages.LoadMode
	BuildEnv   []string
	BuildFlags []string
}

// Parse package parser default implementation
func (p ParserXTool) Parse(ctx context.Context, pkgreg *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Fset:       fset,
		Context:    ctx,
		Dir:        p.AbsDir,
		Mode:       p.LoadMode,
		Env:        p.BuildEnv,
		BuildFlags: p.BuildFlags,
		Tests:      true,
	}
	pkgs, err := packages.Load(cfg, p.Patterns...)
	if err != nil {
		return nil, nil, err
	}
	tpkgs := []*types.Package{}
	for _, pkg := range pkgs {
		if pkgreg.MatchString(pkg.Name) {
			tpkgs = append(tpkgs, pkg.Types)
		}
	}
	return tpkgs, fset, nil
}

// ParserMock defines mock implementation of package parser abstraction
type ParserMock struct {
	pkgs []*types.Package
	fset *token.FileSet
}

// Parse package parser mock implementation
func (p ParserMock) Parse(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	return p.pkgs, p.fset, nil
}

// ParserNil defines package parser nil implementation of parser abstraction
type ParserNil struct{}

// Parse package parser nil implementation
func (ParserNil) Parse(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	return nil, nil, nil
}

// ParserErr defines package parser error implementation of parser abstraction
type ParserErr string

// Parse package parser error implementation
func (p ParserErr) Parse(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	return nil, nil, errors.New(string(p))
}
