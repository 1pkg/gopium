package strategies

import (
	"context"

	"1pkg/gopium"
)

// list of cache presets
var (
	cachel1 = cache{line: 1}
	cachel2 = cache{line: 2}
	cachel3 = cache{line: 3}
)

// cache defines strategy implementation
// that fits structure into cpu cache line
// by adding bottom rounding cpu cache padding
type cache struct {
	curator gopium.Curator
	line    uint
}

// Curator erich cache strategy with curator instance
func (stg cache) Curator(curator gopium.Curator) cache {
	stg.curator = curator
	return stg
}

// Apply cache implementation
func (stg cache) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o
	// calculate size of whole structure
	var size int64
	for _, f := range r.Fields {
		size += f.Size
	}
	// check if cache line size is valid
	if cachel := stg.curator.SysCache(stg.line); cachel > 0 {
		// get number of padding bytes
		// to fill cpu cache line
		// if padding is valid append it
		if pad := size % cachel; pad > 0 {
			pad = cachel - pad
			r.Fields = append(r.Fields, gopium.PadField(pad))
		}
	}
	return r, ctx.Err()
}
