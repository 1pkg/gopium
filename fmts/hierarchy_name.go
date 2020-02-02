package fmts

import (
	"fmt"
	"go/token"
)

// HierarchyName defines abstraction for
// building hierarchy name from flat name, token.FileSet and token.Pos
type HierarchyName func(string, *token.FileSet, token.Pos) string

// FullName defines concating flat name with full hierarchy location HierarchyName implementation
func FullName(name string, fset *token.FileSet, pos token.Pos) string {
	f := fset.File(pos)
	return fmt.Sprintf("%s:L%d %s", f.Name(), f.Line(pos), name)
}

// FlatName defines just flat name returning HierarchyName implementation
func FlatName(name string, fset *token.FileSet, pos token.Pos) string {
	return name
}
