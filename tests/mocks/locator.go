package mocks

import (
	"go/token"

	"1pkg/gopium"
)

// Pos defines mock pos
// data transfer object
type Pos struct {
	ID  string
	Loc string
}

// Locator defines mock locator implementation
type Locator struct {
	Poses map[token.Pos]Pos
}

// ID mock implementation
func (l Locator) ID(pos token.Pos) string {
	// check if we have it in vals
	if t, ok := l.Poses[pos]; ok {
		return t.ID
	}
	// otherwise return default val
	return ""
}

// Loc mock implementation
func (l Locator) Loc(pos token.Pos) string {
	// check if we have it in vals
	if t, ok := l.Poses[pos]; ok {
		return t.Loc
	}
	// otherwise return default val
	return ""
}

// Locator mock implementation
func (l Locator) Locator(string) (gopium.Locator, bool) {
	return l, true
}

// Fset mock implementation
func (l Locator) Fset(string, *token.FileSet) (*token.FileSet, bool) {
	return token.NewFileSet(), true
}

// Root mock implementation
func (l Locator) Root() *token.FileSet {
	return token.NewFileSet()
}
