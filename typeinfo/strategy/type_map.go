package strategy

import (
	"context"
	"go/token"
	"go/types"

	"1pkg/gopium/typeinfo"
)

// typeMap defines strategy type map implementation
// that goes through all structure fields
// extracts typeinfo for each field
// and put it to the typeinfo.TypeInfo map
type typeMap struct {
	e typeinfo.Extractor
	r map[string]typeinfo.TypeInfo
}

// Execute strategy type map implementation
func (tm *typeMap) Execute(ctx context.Context, nm string, st *types.Struct, fset *token.FileSet) error {
	// get number of struct fields
	nf := st.NumFields()
	// refresh result map
	tm.r = make(map[string]typeinfo.TypeInfo, nf)
	for i := 0; i < nf; i++ {
		// get typeinfo
		field := st.Field(i)
		ti := tm.e.Extract(field.Type())
		// put it to the map
		tm.r[field.Name()] = ti
	}
	return nil
}
