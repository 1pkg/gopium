package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// pad defines strategy implementation
// that align all strucutre field
// to sys or max sys padding
// by adding paddings accordingly to system aligments
type pad struct {
	c   gopium.Curator
	sys bool // should max sys padding be used
}

// Apply pad implementation
func (stg pad) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// setup resulted fields list
	var offset, alignment int64 = 0, stg.c.SysAlign()
	fields := make([]gopium.Field, 0, len(r.Fields))
	// go through all fields
	for _, f := range r.Fields {
		// if we wanna use
		// non max system align
		if !stg.sys {
			alignment = f.Align
		}
		// calculate align with padding
		alpad := align(offset, alignment)
		// if padding greater that zero
		// append [pad]byte padding
		if pad := alpad - offset; pad > 0 {
			fields = append(fields, gopium.Field{
				Name:  "_",
				Type:  fmt.Sprintf("[%d]byte", pad),
				Size:  pad,
				Align: 1, // fixed number for byte
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
// note: copied from `go/types/sizes.go`
func align(x, a int64) int64 {
	y := x + a - 1
	return y - y%a
}
