package astutil

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
)

// Apply defines abstraction for
// applying custom action on original ast.Package
// with gopium.Struct map
type Apply func(
	context.Context,
	*ast.Package,
	gopium.Locator,
	collections.Hierarchic,
) (*ast.Package, error)

// Sync implements Apply and combines:
// - sync with fmtio.FSPT helper
// - filter helper
// - note helper
var Sync = combine(
	sync(fmtio.FSPT),
	filter,
	note,
)
