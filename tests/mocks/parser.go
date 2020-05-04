package mocks

import (
	"context"
	"go/ast"
	"go/types"

	"1pkg/gopium"
)

// Parser defines mock parser implementation
type Parser struct {
	Types *types.Package
	Ast   *ast.Package
	Terr  error
	Aerr  error
}

// ParseTypes mock implementation
func (p Parser) ParseTypes(context.Context) (*types.Package, gopium.Locator, error) {
	return p.Types, Locator{}, p.Terr
}

// ParseAst mock implementation
func (p Parser) ParseAst(context.Context) (*ast.Package, gopium.Locator, error) {
	return p.Ast, Locator{}, p.Aerr
}
