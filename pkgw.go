package main

import (
	"context"
	"fmt"
	"go/token"
	"go/types"
	"regexp"
)

// Walker is interface that describes hierarchical walker that
// applies some strategy on ast.StructType
type Walker interface {
	Visit(context.Context, *regexp.Regexp, Strategy)
}

// Pkgw defines package walker struct that is capable of
// walking through all package's structs and apply strategy on them
type Pkgw struct {
	pkg  *types.Package
	fset *token.FileSet
}

// NewPackageWalker creates instance of Pkgw
func NewPackageWalker(ctx context.Context, spkg string, pkgp Pkgp) (*Pkgw, error) {
	// use parser to collect types, fileset and err
	pkg, fset, err := pkgp(ctx, spkg)
	if err != nil {
		return nil, err
	}
	if pkg == nil || fset == nil {
		return nil, fmt.Errorf("package `%s` wasn't found", spkg)
	}
	return &Pkgw{fset: fset, pkg: pkg}, nil
}

// Visit is Pkgw implementation of Walker Visit
// it goes through all struct decls inside the package
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) Visit(ctx context.Context, reg *regexp.Regexp, sg Strategy) {
	fset := pkgw.fset
	s := pkgw.pkg.Scope()
	// go through all names inside the package scope
	for _, name := range s.Names() {
		// check if object name doesn't matches regexp
		if !reg.MatchString(name) {
			continue
		}
		// in case it does and onject is struct
		// then apply strategy
		if st, ok := s.Lookup(name).Type().Underlying().(*types.Struct); ok {
			sg(ctx, st, fset)
		}
	}
}
