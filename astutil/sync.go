package astutil

import (
	"context"
	"go/ast"

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
		sts map[string]gopium.Struct,
	) (*ast.Package, error) {
		// just reuse inner walk helper
		// and apply format to ast
		return walk(ctx, pkg, loc, sts, wact(sta))
	}
}
