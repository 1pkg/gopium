package astext

import (
	"context"
	"errors"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/collections"

	"golang.org/x/tools/go/ast/astutil"
)

// walk defines function that
// helpes walk through ast node
// with comparator function and
// apply some custom action on them
type Walk func(context.Context, ast.Node, Action, Comparator) (ast.Node, error)

// Action defines action
// for ast walk abstraction
// that applies custom action
// on ast type spec node
type Action func(*ast.TypeSpec, gopium.Struct) error

// Comparator defines comparator
// for ast walk abstraction
// that checks if ast type spec node
// needs to be visitted and returns
// relevant gopium struct and flags
type Comparator func(*ast.TypeSpec) (gopium.Struct, bool, bool)

// WalkSt walks through ast struct type
// nodes with comparator function and
// apply some custom action on them
func WalkSt(ctx context.Context, node ast.Node, act Action, cmp Comparator) (ast.Node, error) {
	// tracks error inside astutil apply
	var err error
	// apply astutil apply to parsed ast package
	// and update structure in ast with action
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
						// check if walk should be stopped
						st, skip, brk := cmp(ts)
						switch {
						case skip:
							return true
						case brk:
							err = errors.New("walk has been stoped")
						default:
							// apply action to ast
							err = act(ts, st)
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

// Cmpid helps to create comparator
// which uses match on structs ids
func Cmpid(loc gopium.Locator, h collections.Hierarchic) Comparator {
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

// Cmploc helps to create comparator
// which uses match on sorted
// struct names in provided loc
func Cmploc(loc gopium.Locator, cat string, h collections.Hierarchic) Comparator {
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

// Cmpnote helps to create comparator
// which adapts comparator impl by adding
// check that structure or any structure's
// field has any notes attached to them
func Cmpnote(cmp Comparator) Comparator {
	return func(ts *ast.TypeSpec) (gopium.Struct, bool, bool) {
		// use underlying comparator func
		st, skip, brk := cmp(ts)
		// check if we should process struct
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
			// notes just skip it
			return st, true, false
		}
		// otherwise return underlying
		// comparator func results
		return st, skip, brk
	}
}
