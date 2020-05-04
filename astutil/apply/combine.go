package apply

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
)

// SFN implements apply and combines:
// - sync with fmtio FSPT helper
// - filter helper
// - note helper
var SFN = combine(
	sync(fmtio.FSPT),
	filter,
	note,
)

// combine helps to pipe several
// ast helpers to single apply func
func combine(funcs ...astutil.Apply) astutil.Apply {
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		hsts collections.Hierarchic,
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
			pkg, err = fun(ctx, pkg, loc, hsts)
			// in case of any error
			// just propagate it
			if err != nil {
				return nil, err
			}
		}
		return pkg, nil
	}
}
