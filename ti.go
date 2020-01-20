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

// GetTypeInfo implements TiExt
// using types sizes impl
func GetTypeInfo(t types.Type) TypeInfo {
	// TODO implement it
	return TypeInfo{}
}
