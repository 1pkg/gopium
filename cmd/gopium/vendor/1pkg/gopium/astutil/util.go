package astutil

import (
	"context"
	"go/ast"
	"go/token"
	"io"

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

// Print defines abstraction for
// ast node printing function to io writer
type Print func(io.Writer, *token.FileSet, ast.Node) error

// Persist defines abstraction for
// persisting ast package to fmtio writer
// by using ast print function
type Persist func(
	context.Context,
	fmtio.Writer,
	Print,
	*ast.Package,
	gopium.Locator,
) error
