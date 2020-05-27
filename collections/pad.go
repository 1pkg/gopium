package collections

import (
	"fmt"
	"math"

	"1pkg/gopium"
)

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
