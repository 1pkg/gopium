package apply

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
)

// sync helps to update ast.Package
// accordingly to gopium struct result
// using custom fmtio ast formatter
func sync(sta fmtio.Ast) astutil.Apply {
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		hsts collections.Hierarchic,
	) (*ast.Package, error) {
		// just reuse inner walk helper
		// and apply format to ast
		return walkPkg(
			ctx,
			pkg,
			compid(loc, hsts),
			wact(sta),
		)
	}
}
