package main

import (
	"context"
	"fmt"
	"go/token"
	"go/types"
	"regexp"
)

// Walker defines hierarchical walker abstraction
// that applies some strategy to tree structures
type Walker interface {
	VisitTop(context.Context, *regexp.Regexp, Strategy) // only top level of the tree
	VisitRec(context.Context, *regexp.Regexp, Strategy) // all levels of the tree
}

// Pkgw defines package walker implementation
// that is capable of walking through all package's structs
// and apply specified strategy to them
type Pkgw struct {
	pkg  *types.Package
	fset *token.FileSet
}

// NewPackageWalker creates instance of Pkgw
func NewPackageWalker(ctx context.Context, pkgnm string, pkgp Pkgp) (*Pkgw, error) {
	// use package parser to collect types, fileset and err
	pkg, fset, err := pkgp(ctx, pkgnm)
	if err != nil {
		return nil, err
	}
	if pkg == nil || fset == nil {
		return nil, fmt.Errorf("package %q wasn't found", pkgnm)
	}
	return &Pkgw{fset: fset, pkg: pkg}, nil
}

// VisitTop implements Pkgw Walker VisitTop method
// it goes through all top level struct decls inside the package
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) VisitTop(ctx context.Context, reg *regexp.Regexp, stg Strategy) {
	scope := pkgw.pkg.Scope()
	pkgw.visit(ctx, reg, scope, stg)
}

// VisitRec implements Pkgw Walker VisitRec method
// it goes through all nested levels struct decls inside the package
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) VisitRec(ctx context.Context, reg *regexp.Regexp, stg Strategy) {
	// rec defines recursive function
	// that goes through all nested scopes
	var rec func(scope *types.Scope)
	rec = func(scope *types.Scope) {
		pkgw.visit(ctx, reg, scope, stg)
		for i := 0; i < scope.NumChildren(); i++ {
			chs := scope.Child(i)
			rec(chs)
		}
	}
	scope := pkgw.pkg.Scope()
	rec(scope)
}

// visit helps to implement Pkgw Walker VisitTop and VisitRec methods
// it goes through all struct decls inside the scope
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) visit(ctx context.Context, reg *regexp.Regexp, scope *types.Scope, stg Strategy) {
	fset := pkgw.fset
	// go through all names inside the package scope
	for _, name := range scope.Names() {
		// check if object name doesn't matches regexp
		if !reg.MatchString(name) {
			continue
		}
		// in case it does and onject is struct
		// then apply strategy
		t := scope.Lookup(name).Type()
		if st, ok := t.Underlying().(*types.Struct); ok {
			// TODO hadle this error
			_ = stg(ctx, name, st, fset)
		}
	}
}
