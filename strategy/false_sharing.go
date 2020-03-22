package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
)

// false_sharing defines system struct padding fields strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then adds sys cpu cache padding for each field
type false_sharing struct {
	m gopium.Maven
	l uint
}

// Apply false_sharing implementation
func (stg false_sharing) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// setup resulted fields list
	cachel := stg.m.SysCache(stg.l)
	fields := make([]gopium.Field, 0, len(r.Fields))
	// go through all fields
	for _, f := range r.Fields {
		fields = append(fields, f)
		// if padding greater that zero
		// append [pad]byte padding
		if alpad := f.Size % cachel; alpad > 0 {
			pad := cachel - alpad
			fields = append(fields, gopium.Field{
				Name:  "_",
				Type:  fmt.Sprintf("[%d]byte", pad),
				Size:  pad,
				Align: 1,
			})
		}
	}
	// update fields list
	r.Fields = fields
	return
}
