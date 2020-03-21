package fmts

import (
	"errors"
	"go/ast"
	"sort"

	"1pkg/gopium"
)

// StructToAst defines abstraction for
// formatting original *ast.TypeSpec with gopium.Struct
type StructToAst func(*ast.TypeSpec, gopium.Struct) error

// FSA implements StructToAst and combines:
// - flatten helper
// - shuffle helper
// - annotate helper
var FSA = combine(
	flatten,
	shuffle,
	annotate,
)

// combine helps to pipe several
// ast helpers to single StructToAst func
func combine(funcs ...StructToAst) StructToAst {
	return func(ts *ast.TypeSpec, st gopium.Struct) error {
		// go through all provided funcs
		for _, fun := range funcs {
			// exec single func
			err := fun(ts, st)
			// in case of any error
			// just propagate it
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// flatten helps to make *ast.StructType
// fields list flat by splitting each
// concatenated fields to separate line
func flatten(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("flatten could only be applied to ast.StructType")
	}
	// prepare result list
	list := make([]*ast.Field, 0, tts.Fields.NumFields())
	// iterate over fields list
	for _, field := range tts.Fields.List {
		// for each concatenated name
		// create separate line
		for _, name := range field.Names {
			// copy current field
			f := *field
			// update names list
			f.Names = []*ast.Ident{name}
			// put it to result list
			list = append(list, &f)
		}
	}
	// update structure fields list
	tts.Fields.List = list
	return nil
}

// ShuffleAst helps to sort fields list
// for original *ast.TypeSpec accordingly resulted gopium.Struct
func shuffle(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("shuffle could only be applied to ast.StructType")
	}
	// err holds inner sorting error
	var err error
	// shuffle fields list
	sort.SliceStable(tts.Fields.List, func(i, j int) bool {
		// in case structure isn't flat save error
		// and keep the same order
		if len(tts.Fields.List[i].Names) != 1 || len(tts.Fields.List[j].Names) != 1 {
			err = errors.New("annotate could only be applied to flatten structures")
			return i < j
		}
		// we can safely pick only first name
		// as structure is flat
		// get ast's i-th structure field
		ni := tts.Fields.List[i].Names[0].Name
		// we can safely pick only first name
		// as structure is flat
		// get ast's j-th structure field
		nj := tts.Fields.List[j].Names[0].Name
		// prepare comparison indexes
		// and search for them in resulted structure
		fi, fj := 0, 0
		for index, field := range st.Fields {
			// in case field name of resulted
			// structure matches either:
			// - ast's i-th structure field
			// - ast's j-th structure field
			// set related comparison index
			switch field.Name {
			case ni:
				fi = index
			case nj:
				fj = index
			}
		}
		// compare comparison indexes
		return fi < fj
	})
	// no error can happen
	return err
}

// annotate helps to sync docs and comments
// between original *ast.TypeSpec and resulted gopium.Struct
func annotate(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("annotate could only be applied to ast.StructType")
	}
	// prepare struct docs list
	sdocs := make([]*ast.Comment, 0, len(st.Doc))
	// in case original structure has doc
	if ts.Doc != nil {
		// prepare struct docs list
		sdocs := make([]*ast.Comment, 0, len(ts.Doc.List)+len(st.Doc))
		// collect all docs from original structure
		for _, d := range ts.Doc.List {
			// in case doc has autogenerated prefix skip it
			if !gopium.Stamped(d.Text) {
				sdocs = append(sdocs, d)
			}
		}
	}
	// collect all docs from resulted structure
	for _, d := range st.Doc {
		sdoc := ast.Comment{Text: d}
		sdocs = append(sdocs, &sdoc)
	}
	// update docs list
	ts.Doc = &ast.CommentGroup{List: sdocs}
	// prepare struct comments list
	scomments := make([]*ast.Comment, 0, len(st.Comment))
	// in case original structure has comment
	if ts.Comment != nil {
		// prepare struct comments list
		scomments := make([]*ast.Comment, 0, len(ts.Comment.List)+len(st.Comment))
		// collect all comments from original structure
		for _, c := range ts.Comment.List {
			// in case comment has autogenerated prefix skip it
			if !gopium.Stamped(c.Text) {
				scomments = append(scomments, c)
			}
		}
	}
	// collect all comments from resulted structure
	for _, c := range st.Comment {
		scomment := ast.Comment{Text: c}
		scomments = append(scomments, &scomment)
	}
	// update comments list
	ts.Comment = &ast.CommentGroup{List: scomments}
	// prepare fields storage for docs list and comments list
	stdocs := make(map[string][]*ast.Comment)
	stcomments := make(map[string][]*ast.Comment)
	// go through all resulted structure fields
	for _, field := range st.Fields {
		// collect all docs from resulted structure
		fdocs := make([]*ast.Comment, 0, len(field.Doc))
		for _, d := range field.Doc {
			fdoc := ast.Comment{Text: d}
			fdocs = append(fdocs, &fdoc)
		}
		// collect all comments from resulted structure
		fcomments := make([]*ast.Comment, 0, len(field.Comment))
		for _, c := range field.Comment {
			fcomment := ast.Comment{Text: c}
			fcomments = append(fcomments, &fcomment)
		}
		// put collected results to storage
		stdocs[field.Name] = fdocs
		stcomments[field.Name] = fcomments
	}
	// go through all original structure fields
	for _, field := range tts.Fields.List {
		var fdocs []*ast.Comment
		// in case original field has doc
		if field.Doc != nil {
			// collect all docs from original structure
			fdocs := make([]*ast.Comment, 0, len(field.Doc.List))
			for _, d := range field.Doc.List {
				// in case doc has autogenerated prefix skip it
				if !gopium.Stamped(d.Text) {
					fdocs = append(fdocs, d)
				}
			}
		}
		var fcomments []*ast.Comment
		// in case original field has comment
		if field.Comment != nil {
			// collect all comments from original structure
			fcomments := make([]*ast.Comment, 0, len(field.Comment.List))
			for _, c := range field.Comment.List {
				// in case comment has autogenerated prefix skip it
				if !gopium.Stamped(c.Text) {
					fcomments = append(fcomments, c)
				}
			}
		}
		// in case structure isn't flat return error
		if len(field.Names) != 1 {
			return errors.New("annotate could only be applied to flatten structures")
		}
		// grab the only field name
		fname := field.Names[0].Name
		// if we have docs in storage
		// append them to collected list
		if stdoc, ok := stdocs[fname]; ok {
			fdocs = append(fdocs, stdoc...)
		}
		// if we have comments in storage
		// append them to collected list
		if stcomment, ok := stcomments[fname]; ok {
			fcomments = append(fcomments, stcomment...)
		}
		// update docs and comments list
		field.Doc = &ast.CommentGroup{List: fdocs}
		field.Comment = &ast.CommentGroup{List: fcomments}
	}
	return nil
}