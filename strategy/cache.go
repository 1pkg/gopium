package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// cache defines strategy implementation
// that fits structure into l cpu cache line
// by adding end resulting cpu cache padding
type cache struct {
	c gopium.Curator
	l uint // cache line num
}

// Apply cache implementation
func (stg cache) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// calculate size of whole structure
	var size int64
	for _, f := range r.Fields {
		size += f.Size
	}
	// get number of padding bytes
	// to fill cpu cache line
	cache := stg.c.SysCache(stg.l)
	if pad := size % cache; pad != 0 {
		pad = cache - pad
		r.Fields = append(r.Fields, gopium.Field{
			Name:  "_",
			Type:  fmt.Sprintf("[%d]byte", pad),
			Size:  pad,
			Align: 1, // fixed number for byte
		})
	}
	return
}
