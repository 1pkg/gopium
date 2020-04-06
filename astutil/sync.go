package astutil

import (
	"context"
	"go/ast"
	"go/token"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
)

// sync helps to update ast.Package
// accordingly to gopium.Struct result
// using custom fmtio.StructToAst formatter
func sync(sta fmtio.StructToAst) Apply {
	// bind func
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		hsts HierarchyStructs,
		fsets map[string]*token.FileSet,
	) (*ast.Package, error) {
		// just reuse inner walk helper
		// and apply format to ast
		return walkPkg(
			ctx,
			pkg,
			hierarchy(loc, hsts),
			wact(sta),
		)
	}
}
