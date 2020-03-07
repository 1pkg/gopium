package gopium

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// IDFunc defines abstraction that helpes
// create unique identifier by token.Pos
type IDFunc func(token.Pos) string

// Locator defines abstraction that helpes
// encapsulate pkgs token.FileSet and provides
// some operations on top of it
type Locator token.FileSet

// Sum calculates sha256 hash hex string
// for specified token.Pos in token.FileSet
func (l *Locator) Sum(p token.Pos) string {
	f := (*token.FileSet)(l).File(p)
	r := fmt.Sprintf("%s/%d", f.Name(), f.Line(p))
	h := sha256.Sum256([]byte(r))
	return hex.EncodeToString(h[:])
}

// Fset just returns token.FileSet back
func (l *Locator) Fset() *token.FileSet {
	return (*token.FileSet)(l)
}

// Parser defines abstraction for
// types packages parsing processor
type TypeParser interface {
	ParseTypes(context.Context) (*types.Package, *Locator, error)
}

// Parser defines abstraction for
// ast packages parsing processor
type AstParser interface {
	ParseAst(context.Context) (*ast.Package, *Locator, error)
}

// Parser defines abstraction for packages parsing processor
type Parser interface {
	TypeParser
	AstParser
}
