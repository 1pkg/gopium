package gopium

import (
	"context"
	"go/ast"
)

// XWalker defines ast walker function abstraction
// that walks through type spec ast nodes with provided
// comparator function and applies some custom action
type XWalker interface {
	Walk(context.Context, ast.Node, XAction, XComparator) (ast.Node, error)
}

// XAction defines xwalker action abstraction
// that applies custom action on ast type spec node
type XAction interface {
	Apply(*ast.TypeSpec, Struct) error
}

// Comparator defines xwalker comparator abstraction
// that checks if ast type spec node needs to be visitted
// and returns relevant gopium struct and existing flag
type XComparator interface {
	Check(*ast.TypeSpec) (Struct, bool)
}
