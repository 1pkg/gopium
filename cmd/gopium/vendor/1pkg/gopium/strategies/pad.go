package strategies

import (
	"context"

	"1pkg/gopium"
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
	r := o
	// preset defaults and check that structure has fields
	var offset, stalign, align int64 = 0, 1, stg.curator.SysAlign()
	if flen := len(r.Fields); flen > 0 {
		// setup resulted fields slice
		fields := make([]gopium.Field, 0, flen)
		// go through all fields
		for _, f := range r.Fields {
			// if we wanna use
			// non max system align
			if !stg.sys {
				align = f.Align
			}
			// update struct align size
			if align > stalign {
				stalign = align
			}
			// check that align size is valid
			if align > 0 {
				// calculate align with padding
				alpad := gopium.Align(offset, align)
				// if padding is valid append it
				if pad := alpad - offset; pad > 0 {
					fields = append(fields, gopium.PadField(pad))
				}
				// increment structure offset
				offset = alpad + f.Size
			}
			fields = append(fields, f)
		}
		// check if struct align size is valid
		// and append final padding to structure
		if stalign > 0 {
			// calculate align with padding
			alpad := gopium.Align(offset, stalign)
			// if padding is valid append it
			if pad := alpad - offset; pad > 0 {
				fields = append(fields, gopium.PadField(pad))
			}
		}
		// update resulted fields
		r.Fields = fields
	}
	return r, ctx.Err()
}
