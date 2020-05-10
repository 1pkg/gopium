package gopium

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"
)

// Locator defines abstraction that helps to
// encapsulate pkgs file set related operations
type Locator interface {
	ID(token.Pos) string
	Loc(token.Pos) string
	Locator(string) (Locator, bool)
	Fset(string, *token.FileSet) (*token.FileSet, bool)
	Root() *token.FileSet
}

// TypeParser defines abstraction for
// types packages parsing processor
type TypeParser interface {
	ParseTypes(context.Context, ...byte) (*types.Package, Locator, error)
}

// AstParser defines abstraction for
// ast packages parsing processor
type AstParser interface {
	ParseAst(context.Context, ...byte) (*ast.Package, Locator, error)
}

// Parser defines abstraction that
// aggregates ast and type parsers abstractions
type Parser interface {
	TypeParser
	AstParser
}
