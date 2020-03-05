package fmts

import (
	"1pkg/gopium"
	"go/ast"
	"sort"
)

// StructToAst defines abstraction for
// formatting original *ast.StructType with gopium.Struct
type StructToAst func(*ast.StructType, gopium.Struct) error

// flatten helps to make *ast.StructType
// fields list flat by splitting each
// concatenated fields to separate line
func flatten(stAst *ast.StructType) {
	// prepare result list
	list := make([]*ast.Field, 0, stAst.Fields.NumFields())
	// iterate over fields list
	for _, field := range stAst.Fields.List {
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
	stAst.Fields.List = list
}

// ShuffleAst defines shuffle only
// StructToAst implementation
func ShuffleAst(stAst *ast.StructType, st gopium.Struct) error {
	// make ast flat
	flatten(stAst)
	// shuffle fields list
	sort.SliceStable(stAst.Fields.List, func(i, j int) bool {
		// we can safely pick only first name
		// as structure is flat
		// get ast's i-th structure field
		ni := stAst.Fields.List[i].Names[0].Name
		// we can safely pick only first name
		// as structure is flat
		// get ast's j-th structure field
		nj := stAst.Fields.List[j].Names[0].Name
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
	return nil
}
