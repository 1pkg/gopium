package main

import (
	"context"
	"go/token"
	"go/types"
)

// Strategy defines action abstraction
// that applies some strategy on types struct
type Strategy func(context.Context, *types.Struct, *token.FileSet)

// PkgsSizeMap goes thorught structure calculates size for each field and put it to map
type PkgsSizeMap map[string]uint

// Execute package size map implementation
func (sm PkgsSizeMap) Execute(
	ctx context.Context,
	st *types.Struct,
	fset *token.FileSet,
	tie TiExt,
) {
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		ti := tie(field.Type())
		sm[field.Name()+" "+ti.Name] = ti.Size
	}
}
