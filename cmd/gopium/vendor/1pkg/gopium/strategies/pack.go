package strategies

import (
	"context"
	"sort"

	"1pkg/gopium"
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
type pack struct{}

// Apply pack implementation
func (stg pack) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o
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
