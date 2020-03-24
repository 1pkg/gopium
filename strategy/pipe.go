package strategy

import (
	"context"

	"1pkg/gopium"
)

// pipe defines strategy implementation
// that pipes together set of strategies
// by applying them one after another
type pipe struct {
	stgs []gopium.Strategy
}

// Apply pipe implementation
func (stg pipe) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// go through all inner strategies
	// and apply them one by one
	for _, stg := range stg.stgs {
		r, err = stg.Apply(ctx, o)
		// in case of any error
		// return immediately
		if err != nil {
			return
		}
		// copy result back to
		// original structure
		o = r
	}
	return
}
