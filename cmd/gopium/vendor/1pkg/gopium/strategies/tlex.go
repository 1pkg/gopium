package strategies

import (
	"context"
	"sort"

	"1pkg/gopium"
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
func (stg tlex) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// then execute lexicographical sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// sort depends on type of ordering
		if stg.asc {
			return r.Fields[i].Type < r.Fields[j].Type
		} else {
			return r.Fields[i].Type > r.Fields[j].Type
		}
	})
	return
}
