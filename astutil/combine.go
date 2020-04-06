package astutil

import (
	"context"
	"go/ast"
	"go/token"

	"1pkg/gopium"
)

// combine helps to pipe several
// ast helpers to single Apply func
func combine(funcs ...Apply) Apply {
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		hsts HierarchyStructs,
		fsets map[string]*token.FileSet,
	) (*ast.Package, error) {
		// tracks error inside loop
		var err error
		// go through all provided funcs
		for _, fun := range funcs {
			// manage context actions
			// in case of cancelation
			// stop execution
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			// exec single func
			pkg, err = fun(ctx, pkg, loc, hsts, fsets)
			// in case of any error
			// just propagate it
			if err != nil {
				return nil, err
			}
		}
		return pkg, nil
	}
}
