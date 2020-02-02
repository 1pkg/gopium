package pkgs

import (
	"context"
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/fmts"
)

// Walker defines packages Walker implementation
// that is capable of walking through all package's structs
// and apply specified strategy to them
type Walker struct {
	hn   fmts.HierarchyName
	pkgs []*types.Package
	fset *token.FileSet
}

// NewWalker creates instance of packages Walker
// and requires regex and packages Parser instance to gather all packages
func NewWalker(ctx context.Context, hn fmts.HierarchyName, regex *regexp.Regexp, parser Parser) (Walker, error) {
	if hn == nil {
		return Walker{}, errors.New("hierarchy name wasn't defined")
	}
	// use packages Parser to collect types and fileset
	pkgs, fset, err := parser(ctx, regex)
	if err != nil {
		return Walker{}, err
	}
	// check parse results
	if len(pkgs) == 0 || fset == nil {
		return Walker{}, fmt.Errorf("packages %q wasn't found", regex)
	}
	return Walker{hn: hn, fset: fset, pkgs: pkgs}, nil
}

// VisitTop implements packages Walker VisitTop method
// it goes through all top level struct decls inside the package
// and applies strategy if struct name matches regex
func (w Walker) VisitTop(ctx context.Context, reg *regexp.Regexp, stg gopium.Strategy) {
	for _, pkg := range w.pkgs {
		sc := pkg.Scope()
		w.visit(ctx, reg, sc, stg)
	}
}

// VisitDeep implements packages Walker VisitDeep method
// it goes through all nested levels struct decls inside the package
// and applies strategy if struct name matches regex
func (w Walker) VisitDeep(ctx context.Context, reg *regexp.Regexp, stg gopium.Strategy) {
	// deep defines recursive function
	// that goes through all nested scopes
	var deep func(scope *types.Scope)
	deep = func(scope *types.Scope) {
		w.visit(ctx, reg, scope, stg)
		for i := 0; i < scope.NumChildren(); i++ {
			chs := scope.Child(i)
			deep(chs)
		}
	}
	for _, pkg := range w.pkgs {
		deep(pkg.Scope())
	}
}

// visit helps to implement packages Walker VisitTop and VisitDeep methods
// it goes through all struct decls inside the scope
// and applies strategy if struct name matches regex
func (w Walker) visit(ctx context.Context, reg *regexp.Regexp, scope *types.Scope, stg gopium.Strategy) {
	fset := w.fset
	// go through all names inside the package scope
	for _, name := range scope.Names() {
		// check if object name doesn't matches regex
		if !reg.MatchString(name) {
			continue
		}
		// in case it does and onbject is
		// a type name for struct
		// then apply strategy to it
		tn := scope.Lookup(name)
		if _, ok := tn.(*types.TypeName); !ok {
			continue
		}
		if st, ok := tn.Type().Underlying().(*types.Struct); ok {
			// build hierarchy name here
			hn := w.hn(name, fset, tn.Pos())
			// error should be aways handled by top level strategy
			// so theoretically we will never have error here
			err := stg(ctx, hn, st, fset)
			if err != nil {
				// theoretically imposible case
				panic(err)
			}
		}
	}
}
