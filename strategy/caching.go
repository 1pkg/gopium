package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
)

// caching defines cpu cache line padding fields strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then adds cpu cache line padding to cpu cache line
type caching struct {
	m gopium.Maven
	l uint
}

// Apply caching implementation
func (stg caching) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// calculate size of whole structure
	var size int64
	for _, f := range r.Fields {
		size += f.Size
	}
	// get number of padding bytes
	// to fill cpu cache line
	cache := stg.m.SysCache(stg.l)
	if pad := size % cache; pad != 0 {
		pad = cache - pad
		r.Fields = append(r.Fields, gopium.Field{
			Name:  "_",
			Type:  fmt.Sprintf("[%d]byte", pad),
			Size:  pad,
			Align: 1,
		})
	}
	return
}
