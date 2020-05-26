package mocks

import (
	"context"
	"encoding/json"
	"go/ast"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// Xbytes defines mock fmtio bytes implementation
type Xbytes struct {
	Err error
}

// Bytes mock implementation
func (fmt Xbytes) Bytes(sts []gopium.Struct) ([]byte, error) {
	// in case we have error
	// return it back
	if fmt.Err != nil {
		return nil, fmt.Err
	}
	// otherwise use json bytes impl
	return json.MarshalIndent(sts, "", "\t")
}

// Xast defines mock ast type spec implementation
type Xast struct {
	Err error
}

// Ast mock implementation
func (fmt Xast) Ast(*ast.TypeSpec, gopium.Struct) error {
	return fmt.Err
}

// Diff defines mock diff implementation
type Xdiff struct {
	Err error
}

// Diff mock implementation
func (fmt Xdiff) Diff(o gopium.Categorized, r gopium.Categorized) ([]byte, error) {
	// in case we have error
	// return it back
	if fmt.Err != nil {
		return nil, fmt.Err
	}
	// otherwise use json bytes impl
	fo, fr := collections.Flat(o.Full()), collections.Flat(r.Full())
	data := [][]gopium.Struct{fo.Sorted(), fr.Sorted()}
	return json.MarshalIndent(data, "", "\t")
}

// Xapply defines mock xapply implementation
type Xapply struct {
	Err error
}

// Apply mock implementation
func (a Xapply) Apply(context.Context, *ast.Package, gopium.Locator, gopium.Categorized) (*ast.Package, error) {
	return nil, a.Err
}
