package strategies

import (
	"context"
	"sort"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// list of length presets
var (
	tlenasc  = tlen{asc: true}
	tlendesc = tlen{asc: false}
)

// nlen defines strategy implementation
// that sorts fields accordingly to their types
// length in ascending or descending order
type tlen struct {
	asc bool
}

// Apply tlen implementation
func (stg tlen) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// then execute len sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// sort depends on type of ordering
		if stg.asc {
			return len(r.Fields[i].Type) < len(r.Fields[j].Type)
		}
		return len(r.Fields[i].Type) > len(r.Fields[j].Type)
	})
	return r, ctx.Err()
}
