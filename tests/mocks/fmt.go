package mocks

import (
	"context"
	"encoding/json"
	"go/ast"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// Bytes defines mock fmtio bytes implementation
type Bytes struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Bytes mock implementation
func (fmt Bytes) Bytes(sts []gopium.Struct) ([]byte, error) {
	// in case we have error
	// return it back
	if fmt.Err != nil {
		return nil, fmt.Err
	}
	// otherwise use json bytes impl
	return json.MarshalIndent(sts, "", "\t")
}

// Ast defines mock ast type spec implementation
type Ast struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Ast mock implementation
func (fmt Ast) Ast(*ast.TypeSpec, gopium.Struct) error {
	return fmt.Err
}

// Diff defines mock diff implementation
type Diff struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Diff mock implementation
func (fmt Diff) Diff(o gopium.Categorized, r gopium.Categorized) ([]byte, error) {
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

// Apply defines mock apply implementation
type Apply struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Apply mock implementation
func (a Apply) Apply(context.Context, *ast.Package, gopium.Locator, gopium.Categorized) (*ast.Package, error) {
	return nil, a.Err
}
