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
func (stg pad) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// setup resulted fields list
	var offset, alignment int64 = 0, stg.curator.SysAlign()
	fields := make([]gopium.Field, 0, len(r.Fields))
	// go through all fields
	for _, f := range r.Fields {
		// if we wanna use
		// non max system align
		if !stg.sys {
			alignment = f.Align
		}
		// calculate align with padding
		alpad := gopium.Align(offset, alignment)
		// if padding not equals zero
		// append padding
		if pad := alpad - offset; pad != 0 {
			fields = append(fields, gopium.PadField(pad))
		}
		// increment structure offset
		offset = alpad + f.Size
		fields = append(fields, f)
	}
	// update fields list
	r.Fields = fields
	return
}
