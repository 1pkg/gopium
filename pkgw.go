package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
)

type Pkgw struct {
	fset *token.FileSet
	pkg  *ast.Package
}

func NewPackageWalker(fpkg string) (*Pkgw, error) {
	fset := token.NewFileSet()
	var all = func(os.FileInfo) bool { return true }
	var mode = parser.Mode(0)
	pkgs, err := parser.ParseDir(fset, fpkg, all, mode)
	if err != nil {
		return nil, err
	}

	pkg, ok := pkgs[fpkg]
	if !ok {
		return nil, fmt.Errorf("package %s wasn't found", fpkg)
	}
	return &Pkgw{fset: fset, pkg: pkg}, nil
}

func (pkgw Pkgw) Visit(reg *regexp.Regexp, apply Apply) {
}
