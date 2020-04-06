package typepkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/token"
	"path/filepath"
)

// Locator defines abstraction that helps
// encapsulate pkgs token.FileSet and provides
// some operations on top of it
type Locator token.FileSet

// ID calculates sha256 hash hex string
// for specified token.Pos in token.FileSet
// note: generated id would be ordered inside same loc
func (l *Locator) ID(p token.Pos) string {
	f := (*token.FileSet)(l).File(p)
	// generate hash sum
	r := fmt.Sprintf("%s/%d", f.Name(), f.Line(p))
	h := sha256.Sum256([]byte(r))
	sum := hex.EncodeToString(h[:])
	// generate ordered id
	return fmt.Sprintf("%d-%s", f.Line(p), sum)
}

// Cat returns filepath base file
// for specified token.Pos in token.FileSet
func (l *Locator) Cat(p token.Pos) string {
	f := (*token.FileSet)(l).File(p)
	return filepath.Base(f.Name())
}

// Loc returns filepath base dir
// for specified token.Pos in token.FileSet
func (l *Locator) Loc(p token.Pos) string {
	f := (*token.FileSet)(l).File(p)
	return filepath.Dir(f.Name())
}

// Fset just returns token.FileSet back
func (l *Locator) Fset() *token.FileSet {
	return (*token.FileSet)(l)
}
