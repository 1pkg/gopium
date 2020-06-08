package gopium

import (
	"context"
	"go/ast"
)

// Bytes defines abstraction for formatting
// gopium flat collection to byte slice
type Bytes func([]Struct) ([]byte, error)

// Ast defines abstraction for
// formatting original ast type spec
// accordingly to gopium struct
type Ast func(*ast.TypeSpec, Struct) error

// Diff defines abstraction for formatting
// gopium collections difference to byte slice
type Diff func(Categorized, Categorized) ([]byte, error)

// Apply defines abstraction for
// formatting original ast package by
// applying custom action accordingly to
// provided categorized collection
type Apply func(context.Context, *ast.Package, Locator, Categorized) (*ast.Package, error)
