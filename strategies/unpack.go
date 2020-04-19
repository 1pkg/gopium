package strategies

import (
	"context"
	"math"

	"1pkg/gopium"
)

// list of unpack presets
var (
	unpck = unpack{}
)

// unpack defines strategy implementation
// that rearranges structure fields
// to obtain inflated memory utilization
// by sorting fields accordingly
// to their aligns and sizes in some order
type unpack struct{}

// Apply unpack implementation
func (stg unpack) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// execute pack strategy
	r, err := pck.Apply(ctx, o)
	// check that struct has some fields
	if err != nil || len(r.Fields) == 0 {
		return o, err
	}
	// slice fields by half ceil
	mid := int(math.Ceil(float64(len(r.Fields)) / 2.0))
	left, right := r.Fields[:mid], r.Fields[mid:]
	r.Fields = make([]gopium.Field, 0, len(r.Fields))
	// combine fields in chess order
	for li, ri := 0, len(right)-1; li < mid; li, ri = li+1, ri-1 {
		if ri >= 0 {
			r.Fields = append(r.Fields, right[ri])
		}
		r.Fields = append(r.Fields, left[li])
	}
	return r, ctx.Err()
}
