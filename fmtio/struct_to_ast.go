package fmtio

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"sort"
	"strconv"

	"1pkg/gopium"
)

// StructToAst defines abstraction for
// formatting original ast.TypeSpec with gopium.Struct
type StructToAst func(*ast.TypeSpec, gopium.Struct) error

// FSPTN implements StructToAst and combines:
// - flatten helper
// - fpadfilter helper
// - shuffle helper
// - padsync helper
// - tagsync helper
// - reindex helper
var FSPT = combine(
	flatten,
	fpadfilter,
	shuffle,
	padsync,
	tagsync,
	reindex,
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
		fni := tts.Fields.List[i].Names[0]
		ni := fni.Name
		// we can safely pick only first name
		// as structure is flat
		// get ast's j-th structure field
		fnj := tts.Fields.List[j].Names[0]
		nj := fnj.Name
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
	fields := make([]*ast.Field, len(st.Fields))
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

// reindex helps to reindex fields local token pos
// for original *ast.TypeSpec, by just incrementing
// pos for each struct field, note this is not
// full compliant ast implementation as we are losing
// absolute pos for all other elements, but it's
// too complex to recalculate all elements pos, so
// we can just recalculate local pos which will lead
// to almost identical result
func reindex(ts *ast.TypeSpec, st gopium.Struct) error {
	// check that we are working with ast.StructType
	tts, ok := ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("reindex could only be applied to ast.StructType")
	}
	// set initial pos to zero inside a structure
	pos := token.Pos(0)
	// go through all structure fields
	for _, field := range tts.Fields.List {
		// in case structure isn't flat return error
		if len(field.Names) != 1 {
			return errors.New("reindex could only be applied to flatten structures")
		}
		// set field to current pos
		field.Names[0].NamePos = pos
		// just increment pos
		pos += token.Pos(1)
	}
	return nil
}
