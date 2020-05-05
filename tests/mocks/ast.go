package mocks

import (
	"go/ast"

	"1pkg/gopium"
)

// Ast defines mock fmtio ast implementation
type Ast struct {
	Err error
}

// Ast mock implementation
func (fmt Ast) Ast(*ast.TypeSpec, gopium.Struct) error {
	return fmt.Err
}
