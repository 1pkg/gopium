package astutil

import (
	"context"
	"go/ast"
	"go/token"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"

	"golang.org/x/tools/go/ast/astutil"
)

// walk ines gopium ast walk implementation
// that walks through ast struct type
// nodes with comparator function and
// apply some custom action on them
func walk(ctx context.Context, node ast.Node, act gopium.Action, cmp gopium.Comparator) (ast.Node, error) {
	// err tracks error inside astutil apply
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
						// and skip irrelevant structs
						if st, ok := cmp.Check(ts); ok {
							// apply action to ast
							err = act.Apply(ts, st)
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

// fmtioast defines gopium ast walk
// action fmtio ast implementation
type fmtioast struct {
	fmt fmtio.Ast
}

// Apply fmtioast implementation
func (fmt fmtioast) Apply(ts *ast.TypeSpec, st gopium.Struct) error {
	return fmt.fmt(ts, st)
}

// bcollect defines gopium ast walk
// action boundaries collector implementation
type bcollect struct {
	bs collections.Boundaries
}

// Apply bcollect implementation
func (b *bcollect) Apply(ts *ast.TypeSpec, st gopium.Struct) error {
	// collect structs boundaries
	tts := ts.Type.(*ast.StructType)
	b.bs = append(b.bs, collections.Boundary{
		First: tts.Fields.Opening,
		Last:  tts.Fields.Closing,
	})
	return nil
}

// pressdoccom defines gopium ast walk
// action press doc to file implementation
// which presses comments from
// gopium structure to ast file
type pressdoccom struct {
	file *ast.File
}

// Apply pressdoc implementation
func (pdc pressdoccom) Apply(ts *ast.TypeSpec, st gopium.Struct) error {
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
		pdc.file.Comments = append(pdc.file.Comments, &ast.CommentGroup{List: stdocs})
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
		pdc.file.Comments = append(pdc.file.Comments, &ast.CommentGroup{List: stcoms})
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
			pdc.file.Comments = append(pdc.file.Comments, &ast.CommentGroup{List: fdocs})
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
			pdc.file.Comments = append(pdc.file.Comments, &ast.CommentGroup{List: fcoms})
		}
	}
	return nil
}

// cmpid defines gopium ast walk
// comparator structs flat ids implementation
type cmpid struct {
	loc gopium.Locator
	sts collections.Flat
}

// Check cmpid implementation
func (cmp *cmpid) Check(ts *ast.TypeSpec) (gopium.Struct, bool) {
	// just check if struct
	// with such id is inside
	id := cmp.loc.ID(ts.Pos())
	st, ok := cmp.sts[id]
	return st, ok
}

// cmpsorted defines gopium ast walk
// comparator structs flat implementation
// which uses match on sorted structs name
type cmpsorted struct {
	sts []gopium.Struct
}

// Check cmpsorted implementation
func (cmp cmpsorted) Check(ts *ast.TypeSpec) (gopium.Struct, bool) {
	// if sorted list is not empty
	if len(cmp.sts) > 0 {
		// check the top sorted element name
		if st := cmp.sts[0]; st.Name == ts.Name.Name {
			cmp.sts = cmp.sts[1:]
			return st, true
		}
	}
	// otherwise skip it
	return gopium.Struct{}, false
}

// cmpnote defines gopium ast walk
// comparator adapter implementation
// which adapts provided comparator by adding
// check that structure or any structure's
// field has any notes attached to them
type cmpnote struct {
	cmp gopium.Comparator
}

// Check cmpnote implementation
func (cmp cmpnote) Check(ts *ast.TypeSpec) (gopium.Struct, bool) {
	// use underlying comparator func
	// check if we should process struct
	if st, ok := cmp.cmp.Check(ts); ok {
		// if struct has any notes
		if len(st.Doc) > 0 || len(st.Comment) > 0 {
			return st, true
		}
		// if any field of struct has any notes
		for _, f := range st.Fields {
			if len(f.Doc) > 0 || len(f.Comment) > 0 {
				return st, true
			}
		}
		// in case struct has no inner
		// notes just skip it
		return st, false
	}
	// otherwise skip it
	return gopium.Struct{}, false
}
