package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

// Walker is interface that describes hierarchical walker that
// applies some strategy on ast.StructType
type Walker interface {
	Visit(reg *regexp.Regexp, strg Strategy)
}

// Pkgw defines package walker struct that is capable of
// walking through all package's structs and apply strategy on them
type Pkgw struct {
	fset *token.FileSet
	pkg  *ast.Package
}

// NewPackageWalker creates instance of Pkgw
func NewPackageWalker(fpkg string, pkgp Pkgp) (*Pkgw, error) {
	// use parser to collect pakages, fileset and err
	pkgs, fset, err := pkgp(fpkg)
	if err != nil {
		return nil, err
	}
	// check if pakages list has desired package
	pkg, ok := pkgs[fpkg]
	if !ok {
		return nil, fmt.Errorf("package `%s` wasn't found", fpkg)
	}
	return &Pkgw{fset: fset, pkg: pkg}, nil
}

// Visit is Pkgw implementation of Walker Visit
// it goes through all struct decls inside the package
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) Visit(reg *regexp.Regexp, sg Strategy) {
	// go through all files inside the package
	for _, file := range pkgw.pkg.Files {
		// go through all declarations inside a file
		for _, decl := range file.Decls {
			// check if decl is gendecl
			gdecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			// go through all gendecl specs
			for _, spec := range gdecl.Specs {
				// check if spec is typespec
				tspec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// check if typespec type is structype
				st, ok := tspec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				// check if struct name matches regexp
				// in case it does then apply strategy
				if reg.MatchString(tspec.Name.Name) {
					sg(st)
				}
			}
		}
	}
}
