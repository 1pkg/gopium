package mocks

import (
	"context"
	"go/ast"
	"go/token"
	"io"
)

// Printer defines mock fmtio ast implementation
type Printer struct {
	Err error
}

// Printer mock implementation
func (p Printer) Printer(context.Context, io.Writer, *token.FileSet, ast.Node) error {
	return p.Err
}
