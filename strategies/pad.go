package strategies

import (
	"context"

	"1pkg/gopium/gopium"
	"1pkg/gopium/collections"
)

// list of pad presets
var (
	padsys  = pad{sys: true}
	padtnat = pad{sys: false}
)

// pad defines strategy implementation
// that explicitly aligns each structure field
// to system or type alignment padding
// by adding missing paddings for each field
type pad struct {
	curator gopium.Curator
	sys     bool
}

// Curator erich pad strategy with curator instance
func (stg pad) Curator(curator gopium.Curator) pad {
	stg.curator = curator
	return stg
}

// Apply pad implementation
func (stg pad) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// prepare fields slice
	if flen := len(r.Fields); flen > 0 {
		// set sys align based on sys flag
		sysaling := stg.curator.SysAlign()
		if !stg.sys {
			sysaling = 0
		}
		// collect all struct fields with pads
		rfields := make([]gopium.Field, 0, flen)
		collections.WalkStruct(r, sysaling, func(pad int64, fields ...gopium.Field) {
			// if pad is vallid append it to fields
			if pad > 0 {
				rfields = append(rfields, collections.PadField(pad))
			}
			// append field to fields
			for _, f := range fields {
				rfields = append(rfields, f)
			}
		})
		// update resulted fields
		r.Fields = rfields
	}
	return r, ctx.Err()
}
