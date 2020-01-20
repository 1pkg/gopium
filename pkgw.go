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

// Pkgw defines package walker struct implementation
// that is capable of walking through all package's structs
// and apply specified strategy to them
type Pkgw struct {
	pkg  *types.Package
	fset *token.FileSet
}

// NewPackageWalker creates instance of Pkgw
func NewPackageWalker(ctx context.Context, spkg string, pkgp Pkgp) (*Pkgw, error) {
	// use package parser to collect types, fileset and err
	pkg, fset, err := pkgp(ctx, spkg)
	if err != nil {
		return nil, err
	}
	if pkg == nil || fset == nil {
		return nil, fmt.Errorf("package `%s` wasn't found", spkg)
	}
	return &Pkgw{fset: fset, pkg: pkg}, nil
}

// VisitTop implements Pkgw Walker VisitTop method
// it goes through all top level struct decls inside the package
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) VisitTop(ctx context.Context, reg *regexp.Regexp, sg Strategy) {
	s := pkgw.pkg.Scope()
	pkgw.visit(ctx, reg, s, sg)
}

// VisitRec implements Pkgw Walker VisitRec method
// it goes through all nested levels struct decls inside the package
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) VisitRec(ctx context.Context, reg *regexp.Regexp, sg Strategy) {
	// rec defines recursive function
	// that goes through all nested scopes
	var rec func(s *types.Scope)
	rec = func(s *types.Scope) {
		pkgw.visit(ctx, reg, s, sg)
		for i := 0; i < s.NumChildren(); i++ {
			chs := s.Child(i)
			rec(chs)
		}
	}
	s := pkgw.pkg.Scope()
	rec(s)
}

// visit helps to implement Pkgw Walker VisitTop and VisitRec methods
// it goes through all struct decls inside the scope
// and applies strategy if struct name matches regexp
func (pkgw Pkgw) visit(ctx context.Context, reg *regexp.Regexp, s *types.Scope, sg Strategy) {
	fset := pkgw.fset
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
