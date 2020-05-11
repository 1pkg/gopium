package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium"
)

// Walk defines mock walk implementation
type Walk struct {
	Err error
}

// Walk mock implementation
func (w Walk) Walk(context.Context, ast.Node, gopium.Action, gopium.Comparator) (ast.Node, error) {
	return nil, w.Err
}
