package main

import (
	"context"
	"errors"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// Pkgp is package parser func abstraction for package parsing process
type Pkgp func(context.Context, string) (*types.Package, *token.FileSet, error)

// DefaultPkgp implements Pkgp abstraction and
// executes packages.Load with cfg to collect typs, fileset and err
type DefaultPkgp struct {
	AbsDir     string
	LoadMode   packages.LoadMode
	BuildEnv   []string
	BuildFlags []string
}

// Parse DefaultPkgp implementation
func (pkgp DefaultPkgp) Parse(ctx context.Context, spkg string) (*types.Package, *token.FileSet, error) {
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

// MockPkgp is mock impl of Pkgp abstraction
type MockPkgp struct {
	pkg  *types.Package
	fset *token.FileSet
}

// Parse MockPkgp implementation
func (pkgp MockPkgp) Parse(context.Context, string) (*types.Package, *token.FileSet, error) {
	return pkgp.pkg, pkgp.fset, nil
}

// NotFoundPkgp is not found impl of Pkgp abstraction
type NotFoundPkgp struct{}

// Parse ErrorPkgp implementation
func (NotFoundPkgp) Parse(context.Context, string) (*types.Package, *token.FileSet, error) {
	return nil, nil, nil
}

// ErrorPkgp is error impl of Pkgp abstraction
type ErrorPkgp string

// Parse ErrorPkgp implementation
func (pkgp ErrorPkgp) Parse(context.Context, string) (*types.Package, *token.FileSet, error) {
	return nil, nil, errors.New(string(pkgp))
}
