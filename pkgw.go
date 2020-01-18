package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
)

// Pkgw defines package walker struct that is capable of
// walking through all package's structs and apply action on them
type Pkgw struct {
	fset *token.FileSet
	pkg  *ast.Package
}

// NewPackageWalker creates instance of Pkgw
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

// Visit is Pkgw implementation of Walker Visit
// it goes through all struct decls inside the package
// and applies action if struct name matches regexp
func (pkgw Pkgw) Visit(reg *regexp.Regexp, apply Apply) {
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
				// in case it does then apply action
				if reg.MatchString(tspec.Name.Name) {
					apply(st)
				}
			}
		}
	}
}
