package strategy

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
// additional sys/cpu cache padding
// by adding one padding before and one padding after
// structure fields list
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
func (stg sep) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// get separator size
	sep := stg.curator.SysAlign()
	// if we wanna use
	// non max system separator
	if !stg.sys {
		sep = stg.curator.SysCache(stg.line)
	}
	// add field before or after
	if stg.top {
		r.Fields = append([]gopium.Field{gopium.Pad(sep)}, r.Fields...)
	} else {
		r.Fields = append(r.Fields, gopium.Pad(sep))
	}
	return
}
