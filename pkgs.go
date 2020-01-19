package main

import (
	"context"
	"go/token"
	"go/types"
)

// Strategy is custom callback type that applies some strategy on types.Struct
type Strategy func(context.Context, *types.Struct, *token.FileSet)

// TypeInfo represents specific type information extracted from types.Type
type TypeInfo struct {
	Name string
	Size uint
}

// GetTypeInfo extracts TypeInfo struct from provided types.Type
func GetTypeInfo(t types.Type) TypeInfo {
	// TODO implement it
	return TypeInfo{}
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
