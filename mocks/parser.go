package mocks

import (
	"context"
	"go/ast"
	"go/types"

	"1pkg/gopium"
)

// ParserMock defines mock parser implementation
type ParserMock struct {
	Types *types.Package
	Ast   *ast.Package
	Err   error
}

// ParseTypes mock implementation
func (p ParserMock) ParseTypes(context.Context) (*types.Package, gopium.Locator, error) {
	return p.Types, LocatorMock{}, p.Err
}

// ParseAst mock implementation
func (p ParserMock) ParseAst(context.Context) (*ast.Package, gopium.Locator, error) {
	return p.Ast, LocatorMock{}, p.Err
}
