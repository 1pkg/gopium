package gopium

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"
)

// Locator defines abstraction that helpes
// encapsulate pkgs token.FileSet related operations
type Locator interface {
	ID(token.Pos) string
	Loc(token.Pos) string
	Fset() *token.FileSet
}

// Parser defines abstraction for
// types packages parsing processor
type TypeParser interface {
	ParseTypes(context.Context) (*types.Package, Locator, error)
}

// Parser defines abstraction for
// ast packages parsing processor
type AstParser interface {
	ParseAst(context.Context) (*ast.Package, Locator, error)
}

// Parser defines abstraction for packages parsing processor
type Parser interface {
	TypeParser
	AstParser
}
