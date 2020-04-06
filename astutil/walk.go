package astutil

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"sort"

	"1pkg/gopium"

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
type wcomp func(ts *ast.TypeSpec) (gopium.Struct, bool)

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
						// and skip irrelevant structs
						if st, ok := wcomp(ts); ok {
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

func hierarchy(loc gopium.Locator, hsts HierarchyStructs) wcomp {
	// convert hierarchy struct to flat map
	sts := make(map[string]gopium.Struct, len(hsts))
	for _, lsts := range hsts {
		for id, st := range lsts {
			sts[id] = st
		}
	}
	// return basic comparator func
	return func(ts *ast.TypeSpec) (gopium.Struct, bool) {
		id := loc.ID(ts.Pos())
		st, ok := sts[id]
		return st, ok
	}
}

func ordered(sts map[string]gopium.Struct) wcomp {
	var ids []string
	var stssorted []gopium.Struct
	for id := range sts {
		ids = append(ids, id)
	}
	sort.SliceStable(ids, func(i, j int) bool {
		var idi, idj int
		var sumi, sumj string
		fmt.Sscanf(ids[i], "%d-%s", &idi, &sumi)
		fmt.Sscanf(ids[j], "%d-%s", &idj, &sumj)
		return idi < idj
	})
	for _, id := range ids {
		stssorted = append(stssorted, sts[id])
	}
	return func(ts *ast.TypeSpec) (gopium.Struct, bool) {
		if len(stssorted) > 0 {
			if st := stssorted[0]; st.Name == ts.Name.Name {
				stssorted = stssorted[1:]
				return st, true
			}
		}
		return gopium.Struct{}, false
	}
}
