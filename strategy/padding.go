package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
)

// padding defines system struct padding fields strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then adds field's paddings accordingly to their system aligment
type padding struct {
	m gopium.Maven
}

// Apply padding implementation
func (stg padding) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// setup resulted fields list
	var offset int64
	fields := make([]gopium.Field, 0, len(r.Fields))
	// go through all fields
	for _, f := range r.Fields {
		// calculate align with padding
		alpad := align(offset, f.Align)
		// if padding greater that zero
		// append [pad]byte padding
		if pad := alpad - offset; pad > 0 {
			fields = append(fields, gopium.Field{
				Name:  "_",
				Type:  fmt.Sprintf("[%d]byte", pad),
				Size:  pad,
				Align: 1,
			})
		}
		// increment structure offset
		offset = alpad + f.Size
		fields = append(fields, f)
	}
	// update fields list
	r.Fields = fields
	return
}

// align returns the smallest y >= x such that y % a == 0.
// copied from `go/types/sizes.go`
func align(x, a int64) int64 {
	y := x + a - 1
	return y - y%a
}
