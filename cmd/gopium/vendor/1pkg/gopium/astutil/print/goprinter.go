package print

import (
	"go/ast"
	"go/printer"
	"go/token"
	"io"

	"1pkg/gopium/astutil"
)

// GoPrinter generates go printer ast print instance
// with specified tabwidth and space mode
func GoPrinter(indent int, tabwidth int, usespace bool) astutil.Print {
	// prepare printer with params
	p := &printer.Config{Indent: indent, Tabwidth: tabwidth}
	if usespace {
		p.Mode = printer.UseSpaces
	}
	return func(w io.Writer, fset *token.FileSet, node ast.Node) error {
		// use printer fprint
		return p.Fprint(w, fset, node)
	}
}
