package strategies

import (
	"context"
	"sort"

	"1pkg/gopium/gopium"
	"1pkg/gopium/collections"
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
	asc bool
}

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
