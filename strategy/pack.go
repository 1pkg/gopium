package strategy

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
// that rearranges structure field list
// for optimal memory utilization
// by sorting fields accordingly
// to their sizes in descending order
type pack struct{}

// Apply pack implementation
func (stg pack) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
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
	return
}
