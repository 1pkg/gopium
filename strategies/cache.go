package strategies

import (
	"context"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
)

// list of cache presets
var (
	cachel1  = cache{line: 1, div: true}
	cachel2  = cache{line: 2, div: true}
	cachel3  = cache{line: 3, div: true}
	fcachel1 = cache{line: 1, div: false}
	fcachel2 = cache{line: 2, div: false}
	fcachel3 = cache{line: 3, div: false}
)

// cache defines strategy implementation
// that fits structure into cpu cache line
// by adding bottom rounding cpu cache padding
type cache struct {
	curator gopium.Curator `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	line    uint           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	div     bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	_       [7]byte        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Curator erich cache strategy with curator instance
func (stg cache) Curator(curator gopium.Curator) cache {
	stg.curator = curator
	return stg
}

// Apply cache implementation
func (stg cache) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// calculate aligned size of structure
	var alsize int64
	collections.WalkStruct(r, 0, func(pad int64, fields ...gopium.Field) {
		// add pad to aligned size
		// only if it's not last pad
		if len(fields) > 0 {
			alsize += pad
		}
		// go through all fields
		for _, f := range fields {
			// add field size aligned sizes
			alsize += f.Size
		}
	})
	// check if cache line size is valid
	if cachel := stg.curator.SysCache(stg.line); cachel > 0 {
		// if fractional cache line is allowed
		if stg.div {
			// find smallest size of fraction for cache line
			if alsize > 0 && cachel > alsize {
				for cachel >= alsize && cachel > 1 {
					cachel /= 2
				}
				cachel *= 2
			}
		}
		// get number of padding bytes
		// to fill cpu cache line
		// if padding is valid append it
		if pad := alsize % cachel; pad > 0 {
			pad = cachel - pad
			r.Fields = append(r.Fields, collections.PadField(pad))
		}
	}
	return r, ctx.Err()
}
