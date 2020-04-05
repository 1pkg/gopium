package mocks

import "go/token"

// LocatorMock defines mock locator implementation
type LocatorMock struct {
	IDVal, LocVal string
}

// ID mock implementation
func (l LocatorMock) ID(token.Pos) string {
	return l.IDVal
}

// Loc mock implementation
func (l LocatorMock) Loc(token.Pos) string {
	return l.LocVal
}

// Fset mock implementation
func (l LocatorMock) Fset() *token.FileSet {
	return token.NewFileSet()
}
