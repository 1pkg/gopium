package mocks

import (
	"context"
	"go/ast"
	"go/token"
	"io"
)

// Printer defines mock ast printer implementation
type Printer struct {
	Err error
}

// Print mock implementation
func (p Printer) Print(context.Context, io.Writer, *token.FileSet, ast.Node) error {
	return p.Err
}
