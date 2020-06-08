package gopium

import (
	"context"
	"go/ast"
)

// Walk defines ast walker function abstraction
// that walks through type spec ast nodes with provided
// comparator function and applies some custom action
type Walk func(context.Context, ast.Node, Visitor, Comparator) (ast.Node, error)

// Visitor defines walker action abstraction
// that applies custom action on ast type spec node
type Visitor interface {
	Visit(*ast.TypeSpec, Struct) error
}

// Comparator defines walker comparator abstraction
// that checks if ast type spec node needs to be visitted
// and returns relevant gopium struct and existing flag
type Comparator interface {
	Check(*ast.TypeSpec) (Struct, bool)
}
