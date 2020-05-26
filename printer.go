package gopium

import (
	"context"
	"go/ast"
	"go/token"
	"io"
)

// Printer defines abstraction for
// ast node printing function to io writer
type Printer interface {
	Print(context.Context, io.Writer, *token.FileSet, ast.Node) error
}
