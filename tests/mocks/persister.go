package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium"
)

// Persister defines mock pesister implementation
type Persister struct {
	Err error
}

// Persist mock implementation
func (p Persister) Persist(context.Context, gopium.Printer, gopium.Writer, gopium.Locator, ast.Node) error {
	return p.Err
}
