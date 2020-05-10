package mocks

import (
	"context"
	"go/ast"
	"go/types"

	"1pkg/gopium"
)

// Parser defines mock parser implementation
type Parser struct {
	Typeserr error
	Asterr   error
}

// ParseTypes mock implementation
func (p Parser) ParseTypes(context.Context, ...byte) (*types.Package, gopium.Locator, error) {
	return types.NewPackage("", ""), Locator{}, p.Typeserr
}

// ParseAst mock implementation
func (p Parser) ParseAst(context.Context, ...byte) (*ast.Package, gopium.Locator, error) {
	return &ast.Package{}, Locator{}, p.Asterr
}
