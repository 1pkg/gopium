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
var DefaultPkgp = func(pkg string) (pkgs map[string]*ast.Package, fset *token.FileSet, err error) {
	fset = token.NewFileSet()
	var all = func(os.FileInfo) bool { return true }
	var mode = parser.Mode(0)
	pkgs, err = parser.ParseDir(fset, pkg, all, mode)
	return
}

// MockPkgp is mock impl of Pkgp abstraction
type MockPkgp map[string]*ast.Package

// Parse MockPkgp is mock impl of Pkgp abstraction
func (pkgp MockPkgp) Parse(string) (map[string]*ast.Package, *token.FileSet, error) {
	return map[string]*ast.Package(pkgp), nil, nil
}

// ErrorPkgp is error impl of Pkgp abstraction
type ErrorPkgp string

// ErrorPkgp Parse is error impl of Pkgp abstraction
func (pkgp ErrorPkgp) Parse(string) (map[string]*ast.Package, *token.FileSet, error) {
	return nil, nil, errors.New(string(pkgp))
}
