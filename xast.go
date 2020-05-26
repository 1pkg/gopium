package gopium

import (
	"context"
	"go/ast"
)

// Xwalker defines ast walker function abstraction
// that walks through type spec ast nodes with provided
// comparator function and applies some custom action
type Xwalker interface {
	Walk(context.Context, ast.Node, Xaction, Xcomparator) (ast.Node, error)
}

// Xaction defines xwalker action abstraction
// that applies custom action on ast type spec node
type Xaction interface {
	Apply(*ast.TypeSpec, Struct) error
}

// Xcomparator defines xwalker comparator abstraction
// that checks if ast type spec node needs to be visitted
// and returns relevant gopium struct and existing flag
type Xcomparator interface {
	Check(*ast.TypeSpec) (Struct, bool)
}
