package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// Pkgp is package parser func abstraction for package parsing process
type Pkgp func(string) (map[string]*ast.Package, *token.FileSet, error)

// DefaultPkgp implements Pkgp abstraction and
// executes parser.ParseDir to collect pakages, fileset and err
type DefaultPkgp struct {
	Filter func(os.FileInfo) bool
	Mode   parser.Mode
}

// Parse DefaultPkgp implementation
func (pkgp DefaultPkgp) Parse(pkg string) (pkgs map[string]*ast.Package, fset *token.FileSet, err error) {
	fset = token.NewFileSet()
	pkgs, err = parser.ParseDir(
		fset,
		pkg,
		pkgp.GetFilter(),
		pkgp.GetMode(),
	)
	return
}

// GetFilter gets def pkgp filter with default fallback
func (pkgp DefaultPkgp) GetFilter() func(os.FileInfo) bool {
	if pkgp.Filter != nil {
		return pkgp.Filter
	}
	return func(os.FileInfo) bool { return true }
}

// GetFilter gets def pkgp mode with default fallback
func (pkgp DefaultPkgp) GetMode() parser.Mode {
	if pkgp.Mode != parser.Mode(0) {
		return pkgp.Mode
	}
	return parser.Mode(0)
}

// MockPkgp is mock impl of Pkgp abstraction
type MockPkgp map[string]*ast.Package

// Parse MockPkgp implementation
func (pkgp MockPkgp) Parse(string) (map[string]*ast.Package, *token.FileSet, error) {
	return map[string]*ast.Package(pkgp), nil, nil
}

// ErrorPkgp is error impl of Pkgp abstraction
type ErrorPkgp string

// Parse ErrorPkgp implementation
func (pkgp ErrorPkgp) Parse(string) (map[string]*ast.Package, *token.FileSet, error) {
	return nil, nil, errors.New(string(pkgp))
}
