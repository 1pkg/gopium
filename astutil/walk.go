package astutil

import (
	"context"
	"errors"
	"go/ast"

	"1pkg/gopium"

	"golang.org/x/tools/go/ast/astutil"
)

// wact defines action
// for ast walk helper
// that applies custom action
// on ast.TypeSpec node
type wact func(ts *ast.TypeSpec, st gopium.Struct) error

// walk helps to walk through ast.Package
// with to gopium.Struct result synchronously
// and apply some custom action on them
// after walk returns result or walk error
func walk(
	ctx context.Context,
	pkg *ast.Package,
	loc gopium.Locator,
	sts map[string]gopium.Struct,
	wact wact,
) (*ast.Package, error) {
	// tracks error inside astutil.Apply
	var err error
	// apply astutil.Apply to parsed ast.Package
	// and update structure in ast
	node := astutil.Apply(pkg, func(c *astutil.Cursor) bool {
		if gendecl, ok := c.Node().(*ast.GenDecl); ok {
			for _, spec := range gendecl.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := ts.Type.(*ast.StructType); ok {
						// manage context actions
						// in case of cancelation
						// stop execution
						select {
						case <-ctx.Done():
							err = ctx.Err()
							return false
						default:
						}
						// calculate id for structure
						// and skip all irrelevant structs
						id := loc.ID(ts.Pos())
						if st, ok := sts[id]; ok {
							// apply action to ast
							err = wact(ts, st)
							// in case we have error
							// break iteration
							return err != nil
						}
					}
				}
			}
		}
		return true
	}, nil)
	// in case we had error
	// in astutil.Apply
	// just return it back
	if err != nil {
		return nil, err
	}
	// check that updated type is correct
	if pkg, ok := node.(*ast.Package); ok {
		return pkg, nil
	}
	// in case updated type isn't expected
	return nil, errors.New("can't update package ast")
}
