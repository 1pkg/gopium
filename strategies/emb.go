package strategies

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// list of emb presets
var (
	embasc  = emb{asc: true}
	embdesc = emb{asc: false}
)

// emb defines strategy implementation
// that sorts fields accordingly to their
// embedded flag in ascending or descending order
type emb struct {
	asc bool
}

// Apply emb implementation
func (stg emb) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o
	// then execute embedded sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		if r.Fields[i].Embedded == r.Fields[j].Embedded {
			return false
		}
		// sort depends on type of ordering
		return (r.Fields[i].Embedded && stg.asc) || !(r.Fields[i].Embedded || stg.asc)
	})
	return r, ctx.Err()
}
