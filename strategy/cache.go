package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of cache presets
var (
	cachel1 = cache{l: 1}
	cachel2 = cache{l: 2}
	cachel3 = cache{l: 3}
)

// cache defines strategy implementation
// that fits structure into l cpu cache line
// by adding end resulting cpu cache padding
type cache struct {
	c gopium.Curator
	l uint // cache line num
}

// C erich cache strategy with curator instance
func (stg cache) C(c gopium.Curator) gopium.Strategy {
	stg.c = c
	return stg
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
