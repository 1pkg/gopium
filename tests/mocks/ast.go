package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium/gopium"
)

// Walk defines mock ast walker implementation
type Walk struct {
	Err error
}

// Walk mock implementation
func (w Walk) Walk(context.Context, ast.Node, gopium.Visitor, gopium.Comparator) (ast.Node, error) {
	return nil, w.Err
}
