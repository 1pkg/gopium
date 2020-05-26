package gopium

import (
	"context"
	"go/ast"
)

// Xbytes defines abstraction for formatting
// gopium flat collection to byte slice
type Xbytes func([]Struct) ([]byte, error)

// Xast defines abstraction for
// formatting original ast type spec
// accordingly to gopium struct
type Xast func(*ast.TypeSpec, Struct) error

// Xdiff defines abstraction for formatting
// gopium collections difference to byte slice
type Xdiff func(Categorized, Categorized) ([]byte, error)

// Xapply defines abstraction for
// formatting original ast package by
// applying custom action accordingly to
// provided categorized collection
type Xapply func(context.Context, *ast.Package, Locator, Categorized) (*ast.Package, error)
