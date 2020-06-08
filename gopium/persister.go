package gopium

import (
	"context"
	"go/ast"
)

// Persister defines abstraction for
// ast node pesister with provided printer
// to provided writer by provided locator
type Persister interface {
	Persist(context.Context, Printer, Writer, Locator, ast.Node) error
}
