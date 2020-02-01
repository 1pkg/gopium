package pkgs

import (
	"context"
	"fmt"
	"go/token"
	"go/types"
	"regexp"

	"1pkg/gopium"
)

// Walker defines package walker implementation
// that is capable of walking through all package's structs
// and apply specified strategy to them
type Walker struct {
	pkgs []*types.Package
	fset *token.FileSet
}

// NewWalker creates instance of package walker
// and uses package parser and regex to gather all packages
func NewWalker(ctx context.Context, pkgreg *regexp.Regexp, parser Parser) (*Walker, error) {
	// use package parser to collect types, fileset and err
	pkgs, fset, err := parser(ctx, pkgreg)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 || fset == nil {
		return nil, fmt.Errorf("packages %q wasn't found", pkgreg)
	}
	return &Walker{fset: fset, pkgs: pkgs}, nil
}

// VisitTop implements package walker VisitTop method
// it goes through all top level struct decls inside the package
// and applies strategy if struct name matches regexp
func (w Walker) VisitTop(ctx context.Context, reg *regexp.Regexp, stg gopium.Strategy) {
	for _, pkg := range w.pkgs {
		sc := pkg.Scope()
		w.visit(ctx, reg, sc, stg)
	}
}

// VisitDeep implements package walker VisitDeep method
// it goes through all nested levels struct decls inside the package
// and applies strategy if struct name matches regexp
func (w Walker) VisitDeep(ctx context.Context, reg *regexp.Regexp, stg gopium.Strategy) {
	// rec defines recursive function
	// that goes through all nested scopes
	var rec func(scope *types.Scope)
	rec = func(scope *types.Scope) {
		w.visit(ctx, reg, scope, stg)
		for i := 0; i < scope.NumChildren(); i++ {
			chs := scope.Child(i)
			rec(chs)
		}
	}
	for _, pkg := range w.pkgs {
		rec(pkg.Scope())
	}
}

// visit helps to implement package walker VisitTop and VisitDeep methods
// it goes through all struct decls inside the scope
// and applies strategy if struct name matches regexp
func (w Walker) visit(ctx context.Context, reg *regexp.Regexp, scope *types.Scope, stg gopium.Strategy) {
	fset := w.fset
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
