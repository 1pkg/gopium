package apply

import (
	"context"
	"errors"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/collections"

	"golang.org/x/tools/go/ast/astutil"
)

// wact defines action
// for ast walk helper
// that applies custom action
// on ast.TypeSpec node
type wact func(ts *ast.TypeSpec, st gopium.Struct) error

// wcomp defines comparator
// for ast walk helper
// that checks if ast.TypeSpec node
// needs to be visitted and returns
// relevant struct and visit flag
type wcomp func(ts *ast.TypeSpec) (gopium.Struct, bool, bool)

// walk helps to walk through ast.Node
// on comparator function and
// apply some custom action on them
// after it returns result or error
func walk(
	ctx context.Context,
	node ast.Node,
	wcomp wcomp,
	wact wact,
) (ast.Node, error) {
	// tracks error inside astutil.Apply
	var err error
	// apply astutil.Apply to parsed ast.Package
	// and update structure in ast
	return astutil.Apply(node, func(c *astutil.Cursor) bool {
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
						// check that structure
						// should be visited
						// skip irrelevant structs
						// check if walk should be broken
						if st, skip, brk := wcomp(ts); skip {
							return true
						} else if brk {
							err = errors.New("walk has been stoped")
						} else {
							// apply action to ast
							err = wact(ts, st)
						}
						// in case we have error
						// break iteration
						return err != nil
					}
				}
			}
		}
		return true
	}, nil), err
}

// walkPkg helps to walk through ast.Package
func walkPkg(
	ctx context.Context,
	pkg *ast.Package,
	wcomp wcomp,
	wact wact,
) (*ast.Package, error) {
	// use underlying walk method
	node, err := walk(ctx, pkg, wcomp, wact)
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

// walkFile helps to walk through ast.File
func walkFile(
	ctx context.Context,
	pkg *ast.File,
	wcomp wcomp,
	wact wact,
) (*ast.File, error) {
	// use underlying walk method
	node, err := walk(ctx, pkg, wcomp, wact)
	// in case we had error
	// in astutil.Apply
	// just return it back
	if err != nil {
		return nil, err
	}
	// check that updated type is correct
	if pkg, ok := node.(*ast.File); ok {
		return pkg, nil
	}
	// in case updated type isn't expected
	return nil, errors.New("can't update file ast")
}

// compid helps to create wcomp
// which uses match on structs ids
func compid(loc gopium.Locator, h collections.Hierarchic) wcomp {
	// build flat collection from hierarchic
	f := h.Flat()
	// return basic comparator func
	return func(ts *ast.TypeSpec) (gopium.Struct, bool, bool) {
		// just check if struct
		// with such id is inside
		id := loc.ID(ts.Pos())
		st, ok := f[id]
		return st, !ok, false
	}
}

// comploc helps to create wcomp
// which uses match on sorted
// struct names in provided loc
func comploc(loc gopium.Locator, cat string, h collections.Hierarchic) wcomp {
	// build sorted collection for loc
	f, ok := h.Cat(cat)
	sorted := f.Sorted()
	// return basic comparator func
	return func(ts *ast.TypeSpec) (gopium.Struct, bool, bool) {
		// if loc exists
		if ok && len(sorted) > 0 {
			// check the top sorted element name
			if st := sorted[0]; st.Name == ts.Name.Name {
				sorted = sorted[1:]
				return st, false, false
			}
		}
		// otherwise break it
		return gopium.Struct{}, true, true
	}
}

// compwnote helps to create wcomp
// which adapts wcomp by adding
// check that structure or any structure field
// has any notes inside
func compwnote(comp wcomp) wcomp {
	return func(ts *ast.TypeSpec) (gopium.Struct, bool, bool) {
		// use underlying comp func
		st, skip, brk := comp(ts)
		// check if we can process struct
		if !brk && !skip {
			// if struct has any notes
			if len(st.Doc) > 0 || len(st.Comment) > 0 {
				return st, false, false
			}
			// if any field of struct has any notes
			for _, f := range st.Fields {
				if len(f.Doc) > 0 || len(f.Comment) > 0 {
					return st, false, false
				}
			}
			// in case struct has no inner
			// notes, just skip it
			return st, true, false
		}
		// otherwise return underlying
		// comp func results
		return st, skip, brk
	}
}
