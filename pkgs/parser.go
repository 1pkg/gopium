package pkgs

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"regexp"

	"golang.org/x/tools/go/packages"
)

// Parser defines abstraction for packages parsing processor
type Parser func(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error)

// ParserXTool defines packages Parser default "golang.org/x/tools/go/packages" implementation
// that uses packages.Load with cfg to collect types and fileset
type ParserXTool struct {
	Patterns   []string
	AbsDir     string
	LoadMode   packages.LoadMode
	BuildEnv   []string
	BuildFlags []string
}

// Parse packages Parser default "golang.org/x/tools/go/packages" implementation
func (p ParserXTool) Parse(ctx context.Context, regex *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
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
		if regex.MatchString(pkg.Name) {
			tpkgs = append(tpkgs, pkg.Types)
		}
	}
	return tpkgs, fset, nil
}

// ParserMock defines packages Parser mock implementation
type ParserMock struct {
	pkgs []*types.Package
	fset *token.FileSet
}

// Parse packages Parser mock implementation
func (p ParserMock) Parse(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	return p.pkgs, p.fset, nil
}

// ParserNil defines packages Parser nil implementation
type ParserNil struct{}

// Parse packages Parser nil implementation
func (ParserNil) Parse(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	return nil, nil, nil
}

// ParserError defines packages Parser error implementation
type ParserError string

// Parse packages Parser error implementation
func (p ParserError) Parse(context.Context, *regexp.Regexp) ([]*types.Package, *token.FileSet, error) {
	return nil, nil, errors.New(string(p))
}
