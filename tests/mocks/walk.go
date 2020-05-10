package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium/astext"
)

// Walk defines mock walk implementation
type Walk struct {
	Err error
}

// Walk mock implementation
func (w Walk) Walk(context.Context, ast.Node, astext.Action, astext.Comparator) (ast.Node, error) {
	return nil, w.Err
}
