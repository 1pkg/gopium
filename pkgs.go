package main

import (
	"context"
	"go/token"
	"go/types"
)

// Strategy is custom callback type that applies some strategy on types.Struct
type Strategy func(context.Context, *types.Struct, *token.FileSet)

// SgBuilder is strategy builder interface that helps to create strategy by name
type SgBuilder interface {
	Build(string) (Strategy, error)
}

// PkgsSizeMap goes thorught structure calculates size for each field and put it to map
type PkgsSizeMap map[string]uint

// Execute is Strategy impl for PkgsSizeMap
func (sm PkgsSizeMap) Execute(ctx context.Context, st *types.Struct, fset *token.FileSet) {
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		ti := GetTypeInfo(field.Type())
		sm[field.Name()+" "+ti.Name] = ti.Size
	}
}
