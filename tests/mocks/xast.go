package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium"
)

// Xwalker defines mock ast xwalker implementation
type Xwalker struct {
	Err error
}

// Walk mock implementation
func (w Xwalker) Walk(context.Context, ast.Node, gopium.Xaction, gopium.Xcomparator) (ast.Node, error) {
	return nil, w.Err
}
