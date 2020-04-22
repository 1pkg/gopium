package apply

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/collections"
	"1pkg/gopium/gfmtio/gfmt"
)

// sync helps to update ast.Package
// accordingly to gopium.Struct result
// using custom gfmt.StructToAst formatter
func sync(sta gfmt.StructToAst) astutil.Apply {
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
