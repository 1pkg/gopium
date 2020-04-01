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
// that fits structure into l cpu cache line
// by adding end resulting cpu cache padding
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
	// if padding not equals zero
	// append padding
	cache := stg.curator.SysCache(stg.line)
	if pad := size % cache; pad != 0 {
		pad = cache - pad
		r.Fields = append(r.Fields, gopium.PadField(pad))
	}
	return
}
