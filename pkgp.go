package main

import (
	"context"
	"errors"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// Pkgp defines abstraction for package parsing processor
type Pkgp func(context.Context, string) (*types.Package, *token.FileSet, error)

// PkgpDef defines package parser default implementation
// that uses packages load with cfg to collect types, fileset and err
type PkgpDef struct {
	AbsDir     string
	LoadMode   packages.LoadMode
	BuildEnv   []string
	BuildFlags []string
}

// Parse package parser default implementation
func (pkgp PkgpDef) Parse(ctx context.Context, spkg string) (*types.Package, *token.FileSet, error) {
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Fset:       fset,
		Context:    ctx,
		Dir:        pkgp.AbsDir,
		Mode:       pkgp.LoadMode,
		Env:        pkgp.BuildEnv,
		BuildFlags: pkgp.BuildFlags,
		Tests:      true,
	}
	pkgs, err := packages.Load(cfg, spkg)
	if err != nil {
		return nil, nil, err
	}
	for _, pkg := range pkgs {
		if pkg.Name == spkg {
			return pkg.Types, fset, nil
		}
	}
	return nil, nil, nil
}

// PkgpMock defines mock implementation of pkgp abstraction
type PkgpMock struct {
	pkg  *types.Package
	fset *token.FileSet
}

// Parse package parser mock implementation
func (pkgp PkgpMock) Parse(context.Context, string) (*types.Package, *token.FileSet, error) {
	return pkgp.pkg, pkgp.fset, nil
}

// PkgpNF defines package parser not found implementation of pkgp abstraction
type PkgpNF struct{}

// Parse package parser not found implementation
func (PkgpNF) Parse(context.Context, string) (*types.Package, *token.FileSet, error) {
	return nil, nil, nil
}

// PkgpErr defines package parser error implementation of pkgp abstraction
type PkgpErr string

// Parse package parser error implementation
func (pkgp PkgpErr) Parse(context.Context, string) (*types.Package, *token.FileSet, error) {
	return nil, nil, errors.New(string(pkgp))
}
