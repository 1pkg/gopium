package strategy

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// list of nlex presets
var (
	nlexasc  = nlex{asc: true}
	nlexdesc = nlex{asc: false}
)

// nlex defines strategy implementation
// that sorts fields accordingly to their names
// in ascending or descending order
type nlex struct {
	asc bool
}

// Apply nlex implementation
func (stg nlex) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// then execute lexicographical sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// sort depends on type of ordering
		if stg.asc {
			return r.Fields[i].Name < r.Fields[j].Name
		} else {
			return r.Fields[i].Name > r.Fields[j].Name
		}
	})
	return
}
