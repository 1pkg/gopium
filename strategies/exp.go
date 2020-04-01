package strategies

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// list of exp presets
var (
	expasc  = exp{asc: true}
	expdesc = exp{asc: false}
)

// exp defines strategy implementation
// that sorts fields accordingly to their export flag
// in ascending or descending order
type exp struct {
	asc bool
}

// Apply nlex implementation
func (stg exp) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// then execute exported sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		if r.Fields[i].Exported == r.Fields[j].Exported {
			return false
		}
		// sort depends on type of ordering
		return r.Fields[i].Exported && !stg.asc
	})
	return
}
