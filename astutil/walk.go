package astutil

import (
	"context"
	"errors"
	"go/ast"
	"go/token"

	"1pkg/gopium"
	"1pkg/gopium/collections"

	"golang.org/x/tools/go/ast/astutil"
)

// wact defines action
// for ast walk helper
// that applies custom action
// on ast type spec node
type wact func(ts *ast.TypeSpec, st gopium.Struct) error

// wcomp defines comparator
// for ast walk helper
// that checks if ast type spec node
// needs to be visitted and returns
// relevant gopium struct and flags
type wcomp func(ts *ast.TypeSpec) (gopium.Struct, bool, bool)

// walk helps to walk through ast node
// with comparator function and
// apply some custom action on them
func walk(ctx context.Context, node ast.Node, wcomp wcomp, wact wact) (ast.Node, error) {
	// tracks error inside astutil apply
	var err error
	// apply astutil apply to parsed ast package
	// and update structure in ast with wact
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
						st, skip, brk := wcomp(ts)
						switch {
						case skip:
							return true
						case brk:
							err = errors.New("walk has been stoped")
						default:
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

// pressdoc helps to create wact
// which presses comments from
// gopium structure to ast file
func pressdoc(file *ast.File) wact {
	return func(ts *ast.TypeSpec, st gopium.Struct) error {
		// prepare struct docs slice
		stdocs := make([]*ast.Comment, 0, len(st.Doc))
		// collect all docs from resulted structure
		for _, doc := range st.Doc {
			// doc position is position of name - len of `type` keyword
			slash := ts.Name.Pos() - token.Pos(6)
			sdoc := ast.Comment{Slash: slash, Text: doc}
			stdocs = append(stdocs, &sdoc)
		}
		// update file comments list
		if len(stdocs) > 0 {
			file.Comments = append(file.Comments, &ast.CommentGroup{List: stdocs})
		}
		// prepare struct comments slice
		stcoms := make([]*ast.Comment, 0, len(st.Comment))
		// collect all comments from resulted structure
		for _, com := range st.Comment {
			// comment position is end of type decl
			slash := ts.Type.End()
			scom := ast.Comment{Slash: slash, Text: com}
			stcoms = append(stcoms, &scom)
		}
		// update file comments list
		if len(stcoms) > 0 {
			file.Comments = append(file.Comments, &ast.CommentGroup{List: stcoms})
		}
		// go through all resulted structure fields
		tts := ts.Type.(*ast.StructType)
		for index, field := range st.Fields {
			// get the field from ast
			astfield := tts.Fields.List[index]
			// collect all docs from resulted structure
			fdocs := make([]*ast.Comment, 0, len(field.Doc))
			for _, doc := range field.Doc {
				// doc position is position of name - 1
				slash := astfield.Pos() - token.Pos(1)
				fdoc := ast.Comment{Slash: slash, Text: doc}
				fdocs = append(fdocs, &fdoc)
			}
			// update file comments list
			if len(fdocs) > 0 {
				file.Comments = append(file.Comments, &ast.CommentGroup{List: fdocs})
			}
			// collect all comments from resulted structure
			fcoms := make([]*ast.Comment, 0, len(field.Comment))
			for _, com := range field.Comment {
				// comment position is end of field type
				slash := astfield.Type.End()
				fcom := ast.Comment{Slash: slash, Text: com}
				fcoms = append(fcoms, &fcom)
			}
			// update file comments list
			if len(fcoms) > 0 {
				file.Comments = append(file.Comments, &ast.CommentGroup{List: fcoms})
			}
		}
		return nil
	}
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
// which adapts wcomp impl by adding
// check that structure or any structure's
// field has any notes attached to them
func compwnote(comp wcomp) wcomp {
	return func(ts *ast.TypeSpec) (gopium.Struct, bool, bool) {
		// use underlying comp func
		st, skip, brk := comp(ts)
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
		// comp func results
		return st, skip, brk
	}
}
