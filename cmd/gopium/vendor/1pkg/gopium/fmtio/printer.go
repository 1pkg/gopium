package fmtio

import (
	"context"
	"go/ast"
	"go/printer"
	"go/token"
	"io"

	"1pkg/gopium"

	"golang.org/x/sync/errgroup"
)

// Printer defines abstraction for
// ast node printing function to io writer
type Printer func(io.Writer, *token.FileSet, ast.Node) error

// Goprint generates go printer ast print instance
// with specified tabwidth and space mode
func Goprint(indent int, tabwidth int, usespace bool) Printer {
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

// Save asynchronously pesists ast package
// one ast file by one ast file
// to fmtio writer by using printer
func (p Printer) Save(w Writer) func(ctx context.Context, pkg *ast.Package, loc gopium.Locator) error {
	return func(ctx context.Context, pkg *ast.Package, loc gopium.Locator) error {
		// create sync error group
		// with cancelation context
		group, gctx := errgroup.WithContext(ctx)
		// go through all files in package
		// and update them to concurently
		for name, file := range pkg.Files {
			// manage context actions
			// in case of cancelation
			// stop execution
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}
			// capture name and file copies
			name := name
			file := file
			// run error group write call
			group.Go(func() error {
				// manage context actions
				// in case of cancelation
				// stop execution and return error
				select {
				case <-gctx.Done():
					return gctx.Err()
				default:
				}
				// generate relevant writer
				writer, err := w(name, name)
				// in case any error happened
				// just return error back
				if err != nil {
					return err
				}
				// grab the latest file fset
				fset, _ := loc.Fset(name, nil)
				// write updated ast file to related os file
				// use original file set to keep format
				// in case any error happened
				// just return error back
				if err := p(writer, fset, file); err != nil {
					return err
				}
				// flush writter result
				// in case any error happened
				// just return error back
				return writer.Close()
			})
		}
		// wait until all writers
		// resolve their jobs and
		return group.Wait()
	}
}
