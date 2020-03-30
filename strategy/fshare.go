package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of fshare presets
var (
	fsharel1 = fshare{line: 1}
	fsharel2 = fshare{line: 2}
	fsharel3 = fshare{line: 3}
)

// fshare defines strategy implementation
// that guards structure from false sharing issue
// by adding cpu cache paddings
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
func (stg fshare) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// setup resulted fields list
	cachel := stg.curator.SysCache(stg.line)
	fields := make([]gopium.Field, 0, len(r.Fields))
	// go through all fields
	for _, f := range r.Fields {
		fields = append(fields, f)
		// if padding greater that zero
		// append [pad]byte padding
		if alpad := f.Size % cachel; alpad > 0 {
			pad := cachel - alpad
			fields = append(fields, gopium.Field{
				Name:  "_",
				Type:  fmt.Sprintf("[%d]byte", pad),
				Size:  pad,
				Align: 1, // fixed number for byte
			})
		}
	}
	// update fields list
	r.Fields = fields
	return
}
