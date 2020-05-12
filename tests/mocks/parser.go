package mocks

import (
	"context"
	"go/ast"
	"go/types"

	"1pkg/gopium"
)

// Parser defines mock parser implementation
type Parser struct {
	Parser   gopium.Parser
	Typeserr error
	Asterr   error
}

// ParseTypes mock implementation
func (p Parser) ParseTypes(ctx context.Context, src ...byte) (*types.Package, gopium.Locator, error) {
	// if parser provided use it
	if p.Parser != nil {
		pkg, loc, _ := p.Parser.ParseTypes(ctx, src...)
		return pkg, loc, p.Typeserr
	}
	return types.NewPackage("", ""), Locator{}, p.Typeserr
}

// ParseAst mock implementation
func (p Parser) ParseAst(ctx context.Context, src ...byte) (*ast.Package, gopium.Locator, error) {
	// if parser provided use it
	if p.Parser != nil {
		pkg, loc, _ := p.Parser.ParseAst(ctx, src...)
		return pkg, loc, p.Asterr
	}
	return &ast.Package{}, Locator{}, p.Asterr
}
