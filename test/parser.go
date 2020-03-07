// +build test

package test

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"

	"1pkg/gopium"
)

// ParserMock defines parser mock implementation
type ParserMock struct {
	typePkg *types.Package
	astPkg  *ast.Package
}

// ParseTypes ParserMock implementation
func (p ParserMock) ParseTypes(context.Context) (*types.Package, *gopium.Locator, error) {
	return p.typePkg, (*gopium.Locator)(token.NewFileSet()), nil
}

// ParseAst ParserMock implementation
func (p ParserMock) ParseAst(context.Context) (*ast.Package, *gopium.Locator, error) {
	return p.astPkg, (*gopium.Locator)(token.NewFileSet()), nil
}

// ParserError defines parser error implementation
type ParserError struct {
	err error
}

// ParseTypes ParserError implementation
func (p ParserError) ParseTypes(context.Context) (*types.Package, *gopium.Locator, error) {
	return nil, nil, p.err
}

// ParseAst ParserError implementation
func (p ParserError) ParseAst(context.Context) (*ast.Package, *gopium.Locator, error) {
	return nil, nil, p.err
}
