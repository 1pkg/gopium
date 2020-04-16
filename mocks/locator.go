package mocks

import (
	"go/token"

	"1pkg/gopium"
)

// LocatorMock defines mock locator implementation
type LocatorMock struct {
	IDVal  string
	LocVal string
}

// ID mock implementation
func (l LocatorMock) ID(token.Pos) string {
	return l.IDVal
}

// Loc mock implementation
func (l LocatorMock) Loc(token.Pos) string {
	return l.LocVal
}

// Locator mock implementation
func (l LocatorMock) Locator(string) (gopium.Locator, bool) {
	return l, true
}

// Fset mock implementation
func (l LocatorMock) Fset(string, *token.FileSet) (*token.FileSet, bool) {
	return token.NewFileSet(), true
}

// Root mock implementation
func (l LocatorMock) Root() *token.FileSet {
	return token.NewFileSet()
}
