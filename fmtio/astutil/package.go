package astutil

import (
	"context"
	"go/ast"

	"1pkg/gopium/gopium"

	"golang.org/x/sync/errgroup"
)

// Package ast asyn pesists package implementation
// which persists one ast file by one ast file
// to fmtio writer by using printer
type Package struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Persist package implementation
func (Package) Persist(
	ctx context.Context,
	p gopium.Printer,
	w gopium.Writer,
	loc gopium.Locator,
	node ast.Node,
) error {
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
	// go through all files in package
	// and update them to concurently
	for name, file := range node.(*ast.Package).Files {
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
			// generate relevant writer
			writer, err := w.Generate(name)
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
			if err := p.Print(gctx, writer, fset, file); err != nil {
				return err
			}
			// flush writer result
			// in case any error happened
			// just return error back
			return writer.Close()
		})
	}
	// wait until all writers
	// resolve their jobs and
	return group.Wait()
}
