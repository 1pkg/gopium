package typepkg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/token"
	"sync"

	"1pkg/gopium"
)

// Locator defines abstraction that helps
// encapsulate pkgs token.FileSets and provides
// some operations on top of it
type Locator struct {
	root  *token.FileSet
	extra map[string]*token.FileSet
	mutex sync.Mutex
}

// NewLocator creates new locator instance
// from provided file set
func NewLocator(fset *token.FileSet) *Locator {
	// init root with defaul fset
	// if nil fset has been provided
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
// note: generated ids are ordered inside same loc
func (l *Locator) ID(p token.Pos) string {
	// check if such file exists
	if f := l.root.File(p); f != nil {
		// generate hash sum
		h := sha256.Sum256([]byte(f.Name()))
		sum := hex.EncodeToString(h[:])
		// generate ordered id
		return fmt.Sprintf("%s:%d", sum, f.Line(p))
	}
	return ""
}

// Loc returns full filepath
// for specified token.Pos in token.FileSet
func (l *Locator) Loc(p token.Pos) string {
	// check if such file exists
	if f := l.root.File(p); f != nil {
		return f.Name()
	}
	return ""
}

// Locator returns child locator if any
func (l *Locator) Locator(loc string) (gopium.Locator, bool) {
	fset, ok := l.Fset(loc, nil)
	return NewLocator(fset), ok
}

// Fset multifunc method that
// either set new fset for location
// or returns child fset if any
func (l *Locator) Fset(loc string, fset *token.FileSet) (*token.FileSet, bool) {
	// lock concurrent map access
	defer l.mutex.Unlock()
	l.mutex.Lock()
	// if fset isn't nil
	if fset == nil {
		// write it to exta
		if fset, ok := l.extra[loc]; ok {
			return fset, true
		}
		return l.root, false
	}
	// otherwise read if from exta
	l.extra[loc] = fset
	return fset, true
}

// Root just returns root token.FileSet back
func (l *Locator) Root() *token.FileSet {
	return l.root
}
