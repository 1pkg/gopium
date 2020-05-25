package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium"
)

// XWalker defines mock ast xwalker implementation
type XWalker struct {
	Err error
}

// Walk mock implementation
func (w XWalker) Walk(context.Context, ast.Node, gopium.XAction, gopium.XComparator) (ast.Node, error) {
	return nil, w.Err
}
