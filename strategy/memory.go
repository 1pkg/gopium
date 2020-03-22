package strategy

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// memory defines strategy implementation
// that rearranges structure field list
// for optimal memory utilization
// by sorting fields accordingly
// to their sizes in descending order
type memory struct{}

// Apply memory implementation
func (stg memory) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// execute memory sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		return r.Fields[j].Size < r.Fields[i].Size
	})
	return
}
