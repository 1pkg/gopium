package strategies

import (
	"context"

	"1pkg/gopium"
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
	// preset defaults and check that structure has fields
	var offset, stalign, falign int64 = 0, 1, stg.curator.SysAlign()
	if flen := len(r.Fields); flen > 0 {
		// setup resulted fields slice
		fields := make([]gopium.Field, 0, flen)
		// go through all fields
		for _, f := range r.Fields {
			// if we wanna use
			// non max system align
			if !stg.sys {
				falign = f.Align
			}
			// update struct align size
			if falign > stalign {
				stalign = falign
			}
			// check that align size is valid
			if falign > 0 {
				// calculate align with padding
				alpad := align(offset, falign)
				// if padding is valid append it
				if pad := alpad - offset; pad > 0 {
					fields = append(fields, collections.PadField(pad))
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
			alpad := align(offset, stalign)
			// if padding is valid append it
			if pad := alpad - offset; pad > 0 {
				fields = append(fields, collections.PadField(pad))
			}
		}
		// update resulted fields
		r.Fields = fields
	}
	return r, ctx.Err()
}

// align returns the smallest y >= x such that y % a == 0.
// note: copied from `go/types/sizes.go`
func align(x int64, a int64) int64 {
	y := x + a - 1
	return y - y%a
}
