package astutil

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
)

// Apply defines abstraction for
// applying custom action on original ast.Package
// with gopium.Struct map
type Apply func(context.Context, *ast.Package, gopium.Locator, map[string]gopium.Struct) (*ast.Package, error)

// Sync implements Apply and combines:
// - sync with fmtio.FSPTN helper
// - filternote helper
var Sync = combine(
	sync(fmtio.FSPTN),
	filternote,
)
