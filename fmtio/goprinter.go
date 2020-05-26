package fmtio

import (
	"context"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
)

// Goprinter implements printer by
// using ast go printer printer
type Goprinter struct {
	cfg *printer.Config
}

// NewGoprinter creates instances of goprinter with configs
func NewGoprinter(indent int, tabwidth int, usespace bool) Goprinter {
	// prepare printer with params
	cfg := &printer.Config{Indent: indent, Tabwidth: tabwidth}
	if usespace {
		cfg.Mode = printer.UseSpaces
	}
	return Goprinter{cfg: cfg}
}

// Print goprinter implementation
func (p Goprinter) Print(ctx context.Context, w io.Writer, fset *token.FileSet, node ast.Node) error {
	// manage context actions
	// in case of cancelation
	// stop execution
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// use printer fprint
	return p.cfg.Fprint(w, fset, node)
}
