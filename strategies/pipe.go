package strategies

import (
	"context"

	"1pkg/gopium"
)

// pipe defines strategy implementation
// that pipes together set of strategies
// by applying them one after another
type pipe []gopium.Strategy

// Apply pipe implementation
func (stgs pipe) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o
	// go through all inner strategies
	// and apply them one by one
	for _, stg := range stgs {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-ctx.Done():
			return o, ctx.Err()
		default:
		}
		rst, err := stg.Apply(ctx, r)
		// in case of any error
		// return immediately
		if err != nil {
			return r, err
		}
		// copy result back to
		// result structure
		r = rst
	}
	return r, ctx.Err()
}
