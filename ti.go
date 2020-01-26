package main

import "go/types"

// TypeInfo defines DO type
// that describes some basic type information
type TypeInfo struct {
	Type string
	Size int64
}

// TiExt defines type info extractor abstraction
// to fetch type info from the provided type
type TiExt func(types.Type) TypeInfo

// GetTi implements TiExt
// uses types sizes implentation
func GetTi(t types.Type) TypeInfo {
	sizes := types.SizesFor("gc", "amd64")
	s := sizes.Sizeof(t)
	return TypeInfo{Type: t.String(), Size: s}
}
