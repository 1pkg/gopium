package strategy

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// list of nlen presets
var (
	nlenasc  = nlen{asc: true}
	nlendesc = nlen{asc: false}
)

// nlen defines strategy implementation
// that sorts fields accordingly to their names length
// in ascending or descending order
type nlen struct {
	asc bool
}

// Apply nlen implementation
func (stg nlen) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// then execute len sorting
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
