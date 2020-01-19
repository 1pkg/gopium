package main

import (
	"go/ast"
)

// Strategy is custom callback type that applies some strategy on ast.StructType
type Strategy func(*ast.StructType)

func GetTypeInfo(t ast.Expr) (string, uint) {
	// TODO implement it
	// Initial approach to do everything via ast is overcomplicated.
	// In case we need to find size of nested structs it's ten times simper to use `go/types` pkg then do it buy hand.
	// So need to revisit package walker and package parser in order to update strategy interface to
	// type Strategy func(name string, tt *types.Type, tast *ast.StructType)
	return "", 0
}

// PkgsSizeMap goes thorught structure
// calculates size for each field and put it to map
type PkgsSizeMap map[string]uint

// Execute is Strategy impl for PkgsSizeMap
func (sm PkgsSizeMap) Execute(st *ast.StructType) {
	for _, field := range st.Fields.List {
		tn, sz := GetTypeInfo(field.Type)
		for _, name := range field.Names {
			sm[name.Name+" "+tn] = sz
		}
	}
}
