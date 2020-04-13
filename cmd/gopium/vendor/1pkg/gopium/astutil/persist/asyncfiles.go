package persist

import (
	"context"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/fmtio"

	"golang.org/x/sync/errgroup"
)

// AsyncFiles async pesists ast package to
// writer by using print function
func AsyncFiles(
	ctx context.Context,
	w fmtio.Writer,
	p astutil.Print,
	pkg *ast.Package,
	loc gopium.Locator,
) error {
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
			// in case any error happened just return error
			// it cancels context automatically
			if err != nil {
				return err
			}
			// grab the latest file fset
			fset, _ := loc.Fset(name, nil)
			// write updated ast file to related os file
			// use original file set to keep format
			// in case any error happened just return error
			// it cancels context automatically
			return p(writer, fset, file)
		})
	}
	// wait until all writers
	// resolve their jobs and
	return group.Wait()
}
