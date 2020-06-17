package strategies

import (
	"context"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of fshare presets
var (
	fsharel1 = fshare{line: 1}
	fsharel2 = fshare{line: 2}
	fsharel3 = fshare{line: 3}
)

// fshare defines strategy implementation
// that guards structure from false sharing
// by adding extra cpu cache line paddings
// for each structure field
type fshare struct {
	curator gopium.Curator `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	line    uint           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_       [8]byte        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Curator erich fshare strategy with curator instance
func (stg fshare) Curator(curator gopium.Curator) fshare {
	stg.curator = curator
	return stg
}

// Apply fshare implementation
func (stg fshare) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// check that struct has fields
	// and cache line size is valid
	if flen, cachel := len(r.Fields), stg.curator.SysCache(stg.line); flen > 0 && cachel > 0 {
		// setup resulted fields slice
		fields := make([]gopium.Field, 0, flen)
		// go through all fields
		for _, f := range r.Fields {
			fields = append(fields, f)
			// if padding size is valid
			if pad := f.Size % cachel; pad > 0 {
				pad = cachel - pad
				fields = append(fields, collections.PadField(pad))
			}
		}
		// update resulted fields
		r.Fields = fields
	}
	return r, ctx.Err()
}
