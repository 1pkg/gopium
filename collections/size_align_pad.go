package collections

import (
	"fmt"
	"math"

	"github.com/1pkg/gopium/gopium"
)

// OnPadFields defines pad fields callback that
// accepts pad and optional list of following structure fields
type OnPadFields func(int64, ...gopium.Field)

// WalkStruct iterates over structure fields with optional
// system align and calls on pad fields callback for all pads and fields
func WalkStruct(st gopium.Struct, sysalign int64, onpad OnPadFields) {
	// preset defaults
	var stalign, falign, offset, pad int64 = 1, sysalign, 0, 0
	// calculate total struct size and align
	for _, f := range st.Fields {
		// if provided field align
		// was invalid use field align instead
		if sysalign == 0 {
			falign = f.Align
		}
		// update struct align size
		// if field align size is bigger
		if falign > stalign {
			stalign = falign
		}
		// check that field align is valid
		if falign > 0 {
			// calculate align with padding
			alpad := Align(offset, falign)
			// then calculate padding
			pad = alpad - offset
			// increment structure offset
			offset = alpad + f.Size
		}
		// call onpad func
		onpad(pad, CopyField(f))
		// reset current pad
		pad = 0
	}
	// check if struct align size is valid
	// and append final padding to structure
	if stalign > 0 {
		// calculate align with padding
		alpad := Align(offset, stalign)
		// then calculate padding
		pad = alpad - offset
		// call onpad func
		onpad(pad)
	}
}

// SizeAlign calculates sturct aligned size and size
// by using walk struct helper
func SizeAlign(st gopium.Struct) (int64, int64) {
	// preset defaults
	var alsize, align int64 = 0, 1
	WalkStruct(st, 0, func(pad int64, fields ...gopium.Field) {
		// add pad to aligned size
		alsize += pad
		// go through fields
		for _, f := range fields {
			// add field size aligned sizes
			alsize += f.Size
			// update struct align size
			// if field align size is bigger
			if f.Align > align {
				align = f.Align
			}
		}
	})
	return alsize, align
}

// PadField defines helper that
// creates pad field with specified size
func PadField(pad int64) gopium.Field {
	pad = int64(math.Max(0, float64(pad)))
	return gopium.Field{
		Name:  "_",
		Type:  fmt.Sprintf("[%d]byte", pad),
		Size:  pad,
		Align: 1,
	}
}

// Align returns the smallest y >= x such that y % a == 0.
// note: copied from `go/types/sizes.go`
func Align(x int64, a int64) int64 {
	y := x + a - 1
	return y - y%a
}
