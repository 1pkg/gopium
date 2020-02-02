package fmts

import (
	"fmt"
	"go/token"
	"strings"
)

// HierarchyName defines abstraction for
// building hierarchy name from flat name, token.FileSet and token.Pos
type HierarchyName func(string, *token.FileSet, token.Pos) string

// Root defines abstraction for root part of HierarchyName
// that should be replaced from final FullName result
type Root string

// FullName defines concating flat name with full hierarchy location HierarchyName implementation
func (r Root) FullName(name string, fset *token.FileSet, pos token.Pos) string {
	f := fset.File(pos)
	return fmt.Sprintf(
		"%s:L%d %s",
		strings.Replace(f.Name(), string(r), ".", 1),
		f.Line(pos),
		name,
	)
}

// FlatName defines just flat name returning HierarchyName implementation
func FlatName(name string, fset *token.FileSet, pos token.Pos) string {
	return name
}
