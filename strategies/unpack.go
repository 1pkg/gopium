package strategies

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// list of unpack presets
var (
	unpck = unpack{}
)

// unpack defines strategy implementation
// that rearranges structure field list
// for inflated memory utilization
// by sorting fields accordingly
// to their aligns and sizes in some order
type unpack struct{}

// Apply unpack implementation
func (stg unpack) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// execute memory sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// combine sorting strategy
		// in chess order
		if i%2 == 0 {
			// first compare aligns of two fields
			// lesser aligmnet means upper position
			if r.Fields[i].Align != r.Fields[j].Align {
				return r.Fields[i].Align < r.Fields[j].Align
			}
			// then compare sizes of two fields
			// lesser size means upper position
			return r.Fields[i].Size < r.Fields[j].Size
		} else {
			// first compare aligns of two fields
			// bigger aligmnet means upper position
			if r.Fields[i].Align != r.Fields[j].Align {
				return r.Fields[i].Align > r.Fields[j].Align
			}
			// then compare sizes of two fields
			// bigger size means upper position
			return r.Fields[i].Size > r.Fields[j].Size
		}
	})
	return
}
