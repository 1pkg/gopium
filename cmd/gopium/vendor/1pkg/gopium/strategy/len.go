package strategy

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// list of length presets
var (
	lenasc  = length{true}
	lendesc = length{false}
)

// length defines strategy implementation
// that sorts fields accordingly to their name lengths
// in ascending or descending order
type length struct {
	asc bool
}

// Apply length implementation
func (stg length) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// then execute length sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// sort depends on type of ordering
		if stg.asc {
			return len(r.Fields[i].Name) < len(r.Fields[j].Name)
		} else {
			return len(r.Fields[i].Name) > len(r.Fields[j].Name)
		}
	})
	return
}
