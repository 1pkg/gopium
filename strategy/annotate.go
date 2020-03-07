package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
)

// annotate defines struct size annotation strategy implementation adapter
// that applies underlying strategy and then adds size comment annotation
// for all structure fields and aggregated annotation to structure
type annotate struct {
	stg gopium.Strategy
}

// Apply annotate implementation
func (stg annotate) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply underlying strategy
	o, r, err = stg.stg.Apply(ctx, name, st)
	// then annotate each field with size comment
	var sum int64
	for i := range r.Fields {
		f := &r.Fields[i]
		size := gopium.Stamp(fmt.Sprintf("%d bytes", f.Size))
		f.Comment = append(f.Comment, size)
		sum += f.Size
	}
	// then annotate whole structure with size comment
	size := gopium.Stamp(fmt.Sprintf("%d bytes", sum))
	r.Comment = append(r.Comment, size)
	return
}
