package main

import "go/types"

// TypeInfo defines DO type
// that describes some basic type information
type TypeInfo struct {
	Name string
	Size uint
}

// TiExt defines type info extractor abstraction
// to fetch type info from the provided type
type TiExt func(types.Type) TypeInfo

// GetTi implements TiExt
// uses types sizes implentation
func GetTi(t types.Type) TypeInfo {
	// TODO implement it
	return TypeInfo{}
}
