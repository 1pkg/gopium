package strategies

import (
	"context"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of fshare presets
var (
	sepsyst = sep{sys: true, top: true}
	sepl1t  = sep{line: 1, top: true}
	sepl2t  = sep{line: 2, top: true}
	sepl3t  = sep{line: 3, top: true}
	sepbt   = sep{top: true}
	sepsysb = sep{sys: true, top: false}
	sepl1b  = sep{line: 1, top: false}
	sepl2b  = sep{line: 2, top: false}
	sepl3b  = sep{line: 3, top: false}
	sepbb   = sep{top: false}
)

// sep defines strategy implementation
// that separates structure with
// extra system or cpu cache alignment padding
// by adding the padding at the top
// or the padding at the bottom
type sep struct {
	curator gopium.Curator `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	line    uint           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	bytes   uint           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	sys     bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	top     bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_       [30]byte       `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// Bytes erich sep strategy with custom bytes
func (stg sep) Bytes(bytes uint) sep {
	stg.bytes = bytes
	return stg
}

// Curator erich sep strategy with curator instance
func (stg sep) Curator(curator gopium.Curator) sep {
	stg.curator = curator
	return stg
}

// Apply sep implementation
func (stg sep) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// get separator size
	sep := stg.curator.SysAlign()
	// set separator based on sys flag
	if !stg.sys {
		sep = stg.curator.SysCache(stg.line)
	}
	// if struct has feilds and separator size or bytes are valid
	if flen := len(r.Fields); flen > 0 && (sep > 0 || stg.bytes > 0) {
		if !stg.sys && stg.line == 0 {
			sep = int64(stg.bytes)
		}
		// add field before or after
		// structure fields list
		if stg.top {
			r.Fields = append([]gopium.Field{collections.PadField(sep)}, r.Fields...)
		} else {
			r.Fields = append(r.Fields, collections.PadField(sep))
		}
	}
	return r, ctx.Err()
}
