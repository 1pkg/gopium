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
func (w Walk) Walk(context.Context, ast.Node, astext.Wcmp, astext.Wact) (ast.Node, error) {
	return nil, w.Err
}
