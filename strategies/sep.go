package strategies

import (
	"context"

	"1pkg/gopium"
)

// list of fshare presets
var (
	sepsyst = sep{sys: true, top: true}
	sepl1t  = sep{line: 1, top: true}
	sepl2t  = sep{line: 2, top: true}
	sepl3t  = sep{line: 3, top: true}
	sepsysb = sep{sys: true, top: false}
	sepl1b  = sep{line: 1, top: false}
	sepl2b  = sep{line: 2, top: false}
	sepl3b  = sep{line: 3, top: false}
)

// sep defines strategy implementation
// that separates structure with
// extra system or cpu cache alignment padding
// by adding the padding at the top
// or the padding at the bottom
type sep struct {
	curator gopium.Curator
	line    uint
	sys     bool
	top     bool
}

// Curator erich sep strategy with curator instance
func (stg sep) Curator(curator gopium.Curator) sep {
	stg.curator = curator
	return stg
}

// Apply sep implementation
func (stg sep) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o.Copy()
	// get separator size
	sep := stg.curator.SysAlign()
	// if we wanna use
	// non max system separator
	if !stg.sys {
		sep = stg.curator.SysCache(stg.line)
	}
	// if struct has feilds and separator size is valid
	if flen := len(r.Fields); flen > 0 && sep > 0 {
		// add field before or after
		// structure fields list
		if stg.top {
			r.Fields = append([]gopium.Field{gopium.PadField(sep)}, r.Fields...)
		} else {
			r.Fields = append(r.Fields, gopium.PadField(sep))
		}
	}
	return r, ctx.Err()
}
