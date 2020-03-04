package pkgs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/token"
)

// Locator defines abstraction that helpes
// encapsulate pkgs token.FileSet and provides
// some operations on top of it
type Locator token.FileSet

// NewLocator just follows token.NewFileSet semantic
func NewLocator() *Locator {
	return (*Locator)(token.NewFileSet())
}

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
