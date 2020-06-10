package strategies

import (
	"context"
	"math"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
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
type unpack struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Apply unpack implementation
func (stg unpack) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// execute pack strategy
	r, err := pck.Apply(ctx, r)
	if err != nil {
		return o, err
	}
	// check that struct has fields
	if flen := len(r.Fields); flen > 0 {
		// slice fields by half ceil
		mid := int(math.Ceil(float64(flen) / 2.0))
		left, right := r.Fields[:mid], r.Fields[mid:]
		r.Fields = make([]gopium.Field, 0, flen)
		// combine fields in chess order
		for li, ri := 0, len(right)-1; li < mid; li, ri = li+1, ri-1 {
			if ri >= 0 {
				r.Fields = append(r.Fields, right[ri])
			}
			r.Fields = append(r.Fields, left[li])
		}
	}
	return r, ctx.Err()
}
