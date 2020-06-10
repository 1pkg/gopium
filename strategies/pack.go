package strategies

import (
	"context"
	"sort"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
)

// list of pack presets
var (
	pck = pack{}
)

// pack defines strategy implementation
// that rearranges structure fields
// to obtain optimal memory utilization
// by sorting fields accordingly
// to their aligns and sizes in some order
type pack struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Apply pack implementation
func (stg pack) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// execute memory sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// first compare aligns of two fields
		// bigger aligmnet means upper position
		if r.Fields[i].Align != r.Fields[j].Align {
			return r.Fields[i].Align > r.Fields[j].Align
		}
		// then compare sizes of two fields
		// bigger size means upper position
		return r.Fields[i].Size > r.Fields[j].Size
	})
	return r, ctx.Err()
}
