package mocks

import (
	"go/ast"
	"go/token"
	"io"
)

// Printer defines mock fmtio ast implementation
type Printer struct {
	Err error
}

// Printer mock implementation
func (p Printer) Printer(io.Writer, *token.FileSet, ast.Node) error {
	return p.Err
}
