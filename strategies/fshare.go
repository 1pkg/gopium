package strategies

import (
	"context"

	"1pkg/gopium"
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
	curator gopium.Curator
	line    uint
}

// Curator erich fshare strategy with curator instance
func (stg fshare) Curator(curator gopium.Curator) fshare {
	stg.curator = curator
	return stg
}

// Apply fshare implementation
func (stg fshare) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o
	// check that struct has fields
	// and cache line size is valid
	if flen, cachel := len(r.Fields), stg.curator.SysCache(stg.line); flen > 0 && cachel > 0 {
		// setup resulted fields slice
		fields := make([]gopium.Field, 0, flen)
		// go through all fields
		for _, f := range r.Fields {
			fields = append(fields, f)
			// if padding not equals zero
			// append padding
			if pad := f.Size % cachel; pad > 0 {
				pad = cachel - pad
				fields = append(fields, gopium.PadField(pad))
			}
		}
		// update resulted fields
		r.Fields = fields
	}
	return r, ctx.Err()
}
