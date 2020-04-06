package astutil

import (
	"context"
	"go/ast"
	"go/token"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
)

type HierarchyStructs map[string]map[string]gopium.Struct

// Apply defines abstraction for
// applying custom action on original ast.Package
// with gopium.Struct map
type Apply func(
	context.Context,
	*ast.Package,
	gopium.Locator,
	HierarchyStructs,
	map[string]*token.FileSet,
) (*ast.Package, error)

// Sync implements Apply and combines:
// - sync with fmtio.FSPT helper
// - filternote helper
var Sync = combine(
	sync(fmtio.FSPT),
	filternote,
	note,
)
