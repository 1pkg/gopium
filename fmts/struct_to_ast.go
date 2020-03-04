package fmts

import (
	"1pkg/gopium"
	"go/ast"
)

// StructToAst defines abstraction for
// formatting gopium.Struct to *ast.TypeSpec
type StructToAst func(gopium.Struct) (*ast.TypeSpec, error)
