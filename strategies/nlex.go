package strategies

import (
	"context"
	"sort"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of nlex presets
var (
	nlexasc  = nlex{asc: true}
	nlexdesc = nlex{asc: false}
)

// nlex defines strategy implementation
// that sorts fields accordingly to their names
// in ascending or descending order
type nlex struct {
	asc bool    `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_   [1]byte `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 2 bytes; struct align: 1 bytes; struct aligned size: 2 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg

// Apply nlex implementation
func (stg nlex) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// then execute lexicographical sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// sort depends on type of ordering
		if stg.asc {
			return r.Fields[i].Name < r.Fields[j].Name
		}
		return r.Fields[i].Name > r.Fields[j].Name
	})
	return r, ctx.Err()
}
