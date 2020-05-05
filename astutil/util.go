package astutil

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// Apply defines abstraction for
// applying custom action
// on original ast package
// accordingly to gopium
// hierarchic collection
type Apply func(
	context.Context,
	*ast.Package,
	gopium.Locator,
	collections.Hierarchic,
) (*ast.Package, error)
