package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// Apply defines mock astutil apply implementation
type Apply struct {
	Err error
}

// Apply mock implementation
func (a Apply) Apply(context.Context, *ast.Package, gopium.Locator, collections.Hierarchic) (*ast.Package, error) {
	return nil, a.Err
}
