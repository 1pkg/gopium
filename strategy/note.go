package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of note presets
var (
	nt = note{}
)

// note defines strategy implementation
// that adds size comment annotation
// for all structure fields
// and aggregated size annotation
// for whole structure
type note struct{}

// Apply note implementation
func (stg note) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// note each field with size comment
	var sum int64
	for i := range r.Fields {
		f := &r.Fields[i]
		size := gopium.Stamp(fmt.Sprintf("%d bytes", f.Size))
		f.Comment = append(f.Comment, size)
		sum += f.Size
	}
	// note whole structure with size comment
	size := gopium.Stamp(fmt.Sprintf("%d bytes", sum))
	r.Comment = append(r.Comment, size)
	return
}
