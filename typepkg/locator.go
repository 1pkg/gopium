package typepkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/token"

	"1pkg/gopium"
)

// Locator defines abstraction that helps
// encapsulate pkgs token.FileSets and provides
// some operations on top of it
type Locator struct {
	root  *token.FileSet
	extra map[string]*token.FileSet
}

// NewLocator creates new locator instance
// from provided file set
func NewLocator(fset *token.FileSet) *Locator {
	if fset == nil {
		fset = token.NewFileSet()
	}
	return &Locator{
		root:  fset,
		extra: make(map[string]*token.FileSet),
	}
}

// ID calculates sha256 hash hex string
// for specified token.Pos in token.FileSet
// note: generated id would be ordered inside same loc
func (l *Locator) ID(p token.Pos) string {
	f := l.root.File(p)
	// generate hash sum
	r := fmt.Sprintf("%s/%d", f.Name(), f.Line(p))
	h := sha256.Sum256([]byte(r))
	sum := hex.EncodeToString(h[:])
	// generate ordered id
	return fmt.Sprintf("%d-%s", f.Line(p), sum)
}

// Loc returns full filepath
// for specified token.Pos in token.FileSet
func (l *Locator) Loc(p token.Pos) string {
	return l.root.File(p).Name()
}

// Locator returns child locator if any
func (l *Locator) Locator(loc token.Pos) (gopium.Locator, bool) {
	fset, ok := l.Fset(loc, nil)
	return NewLocator(fset), ok
}

// Fset multifunc method that
// either set new fset for location
// or returns child fset if any
func (l *Locator) Fset(pos token.Pos, fset *token.FileSet) (*token.FileSet, bool) {
	if fset == nil {
		if fset, ok := l.extra[l.Loc(pos)]; ok {
			return fset, true
		}
		return l.root, false
	}
	l.extra[l.Loc(pos)] = fset
	return fset, true
}

// Root just returns root token.FileSet back
func (l *Locator) Root() *token.FileSet {
	return l.root
}
