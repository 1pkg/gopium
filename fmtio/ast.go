package fmtio

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strconv"

	"github.com/1pkg/gopium/gopium"
)

// FSPT implements ast and combines:
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
// ast helpers to single ast func
func combine(funcs ...gopium.Ast) gopium.Ast {
	return func(ts *ast.TypeSpec, st gopium.Struct) error {
		// check that we are working with ast struct type
		if _, ok := ts.Type.(*ast.StructType); !ok {
			return fmt.Errorf("type %q is not valid structure", ts.Name.Name)
		}
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

// flatten helps to make ast struct type
// fields list flat by splitting each
// concatenated fields to separate line
func flatten(ts *ast.TypeSpec, st gopium.Struct) error {
	// prepare result slice
	tts := ts.Type.(*ast.StructType)
	fields := make([]*ast.Field, 0, tts.Fields.NumFields())
	// iterate over fields list
	for _, field := range tts.Fields.List {
		// for each concatenated name
		// create separate line
		for _, name := range field.Names {
			// copy current field
			f := *field
			// update names slice
			f.Names = []*ast.Ident{name}
			// put it to result slice
			fields = append(fields, &f)
		}
		// embedded fields should
		// be still collected
		if len(field.Names) == 0 {
			fields = append(fields, field)
		}
	}
	// update structure fields list
	tts.Fields.List = fields
	return nil
}

// fpadfilter helps to filter fields and pads
// from fields list for original ast type spec
// accordingly to result gopium struct
func fpadfilter(ts *ast.TypeSpec, st gopium.Struct) error {
	// collect unique fields
	tts := ts.Type.(*ast.StructType)
	fields := make(map[string]struct{}, len(st.Fields))
	for _, f := range st.Fields {
		fields[f.Name] = struct{}{}
	}
	// prepare resulted fields slice
	nfields := make([]*ast.Field, 0, len(tts.Fields.List))
	// go through original ast fields list
	for _, f := range tts.Fields.List {
		// start with non embedded fields
		if len(f.Names) == 1 {
			// if pad field was detected
			// filter it out
			if f.Names[0].Name == "_" {
				continue
			}
			// if field isn't inside
			// filter it out
			if _, ok := fields[f.Names[0].Name]; !ok {
				continue
			}
			// otherwise collect field
			nfields = append(nfields, f)
		}
		// embedded fields should
		// be still collected
		if len(f.Names) == 0 {
			nfields = append(nfields, f)
		}
	}
	// update original ast fields list
	tts.Fields.List = nfields
	return nil
}

// shuffle helps to sort fields list
// for ast type spec accordingly to result struct
func shuffle(ts *ast.TypeSpec, st gopium.Struct) error {
	// collect fields indexes
	tts := ts.Type.(*ast.StructType)
	fields := make(map[string]int, len(st.Fields))
	// in case of embedded fields
	// use types to uniquely discern them
	for i, f := range st.Fields {
		if f.Name == "" {
			fields[f.Type] = i
			continue
		}
		fields[f.Name] = i
	}
	// shuffle fields list
	sort.SliceStable(tts.Fields.List, func(i, j int) bool {
		// we can safely pick only first name
		// for flat structure non embedded
		// ast's i-th and j-th fields
		// in case fields are embedded
		// use type instead if possible
		var ni, nj string
		if fni := tts.Fields.List[i]; len(fni.Names) == 1 {
			ni = fni.Names[0].Name
		} else if len(fni.Names) == 0 {
			if it, ok := fni.Type.(*ast.Ident); ok {
				ni = it.Name
			}
		}
		if fnj := tts.Fields.List[j]; len(fnj.Names) == 1 {
			nj = fnj.Names[0].Name
		} else if len(fnj.Names) == 0 {
			if it, ok := fnj.Type.(*ast.Ident); ok {
				nj = it.Name
			}
		}
		// prepare comparison indexes
		// and search for them in resulted structure
		// in case field name of resulted
		// structure matches either:
		// - ast's i-th structure field
		// - ast's j-th structure field
		// set related comparison index
		fi, fj := 0, 0
		if index, ok := fields[ni]; ok {
			fi = index
		}
		if index, ok := fields[nj]; ok {
			fj = index
		}
		// compare comparison indexes
		return fi < fj
	})
	return nil
}

// padsync helps to sync fields padding list
// for ast type spec accordingly to result struct
func padsync(ts *ast.TypeSpec, st gopium.Struct) error {
	// prepare resulted fields slice
	tts := ts.Type.(*ast.StructType)
	fields := make([]*ast.Field, len(st.Fields))
	copy(fields, tts.Fields.List)
	for index, f := range st.Fields {
		// skip non pad fields
		if f.Name != "_" {
			continue
		}
		// transform size to string format
		// and add pad field to struct
		size := strconv.Itoa(int(f.Size))
		field := &ast.Field{
			Names: []*ast.Ident{
				{
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
// ast type spec and result struct
func tagsync(ts *ast.TypeSpec, st gopium.Struct) error {
	// go through all original structure fields
	tts := ts.Type.(*ast.StructType)
	for index, field := range tts.Fields.List {
		// check if field tag exists
		f := st.Fields[index]
		if f.Tag != "" {
			// update ast tag
			field.Tag = &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`%s`", f.Tag),
			}
		}
	}
	return nil
}

// reindex helps to reindex fields local token pos
// for original ast type spec, by just incrementing
// pos for each struct field,
//
// note this is not full compliant ast implementation
// as we are losing absolute pos for all other elements,
//
// but it's too complex to recalculate all elements pos,
// so we can just recalculate local pos which leads to
// almost identical result
func reindex(ts *ast.TypeSpec, st gopium.Struct) error {
	// set initial pos to zero inside a structure
	pos := token.Pos(0)
	// go through all structure fields
	tts := ts.Type.(*ast.StructType)
	for _, field := range tts.Fields.List {
		// in case field isn't flat skip it
		if len(field.Names) == 1 {
			// set field to current pos
			field.Names[0].NamePos = pos
			// just increment pos
			pos += token.Pos(1)
		}
	}
	return nil
}
