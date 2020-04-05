package fmtio

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"1pkg/gopium"
)

// StructToAst defines abstraction for
// formatting original *ast.TypeSpec with gopium.Struct
type StructToAst func(*ast.TypeSpec, gopium.Struct) error

// FSPTN implements StructToAst and combines:
// - flatten helper
// - padfilter helper
// - shuffle helper
// - padsync helper
// - tagsync helper
// - notesync helper
var FSPTN = combine(
	flatten,
	fpadfilter,
	shuffle,
	padsync,
	tagsync,
	notesync,
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

// fpadfilter helps to filter fields and pads
// from fields list for original *ast.TypeSpec
// accordingly to result gopium.Struct
func fpadfilter(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("fpadfilter could only be applied to ast.StructType")
	}
	// collect all unique fields list
	stfields := make(map[string]struct{}, len(st.Fields))
	for _, f := range st.Fields {
		stfields[f.Name] = struct{}{}
	}
	// prepare resulted fields list
	fields := make([]*ast.Field, 0, len(tts.Fields.List))
	// go through original ast struct
	for _, f := range tts.Fields.List {
		// in case structure isn't flat return error
		if len(f.Names) != 1 {
			return errors.New("fpadfilter could only be applied to flatten structures")
		}
		// if pad field was detected
		// filter it out
		if f.Names[0].Name == "_" {
			continue
		}
		// if field isn't in result list
		// filter it out
		if _, ok := stfields[f.Names[0].Name]; !ok {
			continue
		}
		// otherwise collect field
		fields = append(fields, f)
	}
	// update original ast fields list
	tts.Fields.List = fields
	return nil
}

// shuffle helps to sort fields list
// for original *ast.TypeSpec accordingly to result gopium.Struct
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
			case "_": // skip paddings
				index--
			}
		}
		// compare comparison indexes
		return fi < fj
	})
	// no error can happen
	return err
}

// padsync helps to sync fields padding list
// for original *ast.TypeSpec accordingly to result gopium.Struct
func padsync(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("padsync could only be applied to ast.StructType")
	}
	// prepare pad type expression regex
	regex := regexp.MustCompile(`\[.*\]byte`)
	// prepare resulted fields list
	fields := make([]*ast.Field, 0, len(tts.Fields.List)+len(st.Fields))
	copy(fields, tts.Fields.List)
	for index, f := range st.Fields {
		// skip non pad fields
		if f.Name != "_" {
			continue
		}
		// in case pad type is unexpected
		// return error
		if !regex.MatchString(f.Type) {
			return fmt.Errorf("padsync unexpected pad type expression %s", f.Type)
		}
		// transform size to string format
		// and add pad field to struct
		// note: don't need to sync docs/comments here
		// as it will be done in annotate
		size := strconv.Itoa(int(f.Size))
		field := &ast.Field{
			Names: []*ast.Ident{
				&ast.Ident{
					Name: "_",
					Obj: &ast.Object{
						Kind: ast.Var,
						Name: "_",
					},
				},
			},
			Type: &ast.ArrayType{
				Len: &ast.BasicLit{
					Kind:  token.INT,
					Value: size,
				},
				Elt: &ast.Ident{
					Name: "byte",
				},
			},
		}
		// shift fields one right
		copy(fields[index+1:], fields[index:])
		// insert pad at index
		fields[index] = field
	}
	// update original ast fields list
	tts.Fields.List = fields
	return nil
}

// tagsync helps to sync field tags between
// original *ast.TypeSpec result gopium.Struct
func tagsync(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("tagsync could only be applied to ast.StructType")
	}
	// prepare struct tags list
	sttags := make(map[string]*ast.BasicLit)
	// go through all resulted structure fields
	for _, field := range st.Fields {
		// put ast tag to the map
		sttags[field.Name] = &ast.BasicLit{
			Kind:  token.STRING,
			Value: field.Tag,
		}
	}
	// go through all original structure fields
	for _, field := range tts.Fields.List {
		// in case structure isn't flat return error
		if len(field.Names) != 1 {
			return errors.New("tagsync could only be applied to flatten structures")
		}
		// grab the only field name
		fname := field.Names[0].Name
		// if we have tag in the map
		// set it as field tag
		if sttag, ok := sttags[fname]; ok {
			field.Tag = sttag
		}
	}
	return nil
}

// notesync helps to sync docs and comments
// between original *ast.TypeSpec and result gopium.Struct
func notesync(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("notesync could only be applied to ast.StructType")
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
			if !strings.Contains(d.Text, gopium.STAMP) {
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
			if !strings.Contains(c.Text, gopium.STAMP) {
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
				if !strings.Contains(d.Text, gopium.STAMP) {
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
				if !strings.Contains(c.Text, gopium.STAMP) {
					fcomments = append(fcomments, c)
				}
			}
		}
		// in case structure isn't flat return error
		if len(field.Names) != 1 {
			return errors.New("notesync could only be applied to flatten structures")
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
