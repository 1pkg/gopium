package astutil

import (
	"context"
	"go/ast"
	"go/token"
	"io"

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

// Print defines abstraction for
// ast node printing function to io writer
type Print func(io.Writer, *token.FileSet, ast.Node) error

// Persist defines abstraction for
// persisting ast package
// with ast print function
type Persist func(context.Context, Print, *ast.Package, gopium.Locator) error
