package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
)

// stamp defines struct stamp strategy implementation adapter
// that applies underlying strategy and then adds doc stamp to structure
type stamp struct {
	stg gopium.Strategy
}

// Apply stamp implementation
func (stg stamp) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply underlying strategy
	o, r, err = stg.stg.Apply(ctx, name, st)
	// create stamp
	stamp := fmt.Sprintf("struct %q has been auto curated", name)
	stamp = gopium.Stamp(stamp)
	// add stamp to structure doc
	r.Doc = append(r.Doc, stamp)
	return
}
