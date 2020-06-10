package strategies

import (
	"context"
	"sort"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
)

// list of tlex presets
var (
	tlexasc  = tlex{asc: true}
	tlexdesc = tlex{asc: false}
)

// tlex defines strategy implementation
// that sorts fields accordingly to their types
// in ascending or descending order
type tlex struct {
	asc bool    `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	_   [1]byte `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
} // struct size: 2 bytes; struct align: 1 bytes; struct aligned size: 2 bytes; - 🌺 gopium @1pkg

// Apply tlex implementation
func (stg tlex) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// then execute lexicographical sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// sort depends on type of ordering
		if stg.asc {
			return r.Fields[i].Type < r.Fields[j].Type
		}
		return r.Fields[i].Type > r.Fields[j].Type
	})
	return r, ctx.Err()
}
