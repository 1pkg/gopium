package mocks

import (
	"go/token"

	"1pkg/gopium"
)

// Locator defines mock locator implementation
type Locator struct {
	IDVal  string
	LocVal string
}

// ID mock implementation
func (l Locator) ID(token.Pos) string {
	return l.IDVal
}

// Loc mock implementation
func (l Locator) Loc(token.Pos) string {
	return l.LocVal
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
