package strategies

import (
	"context"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// pipe defines strategy implementation
// that pipes together set of strategies
// by applying them one after another
type pipe []gopium.Strategy

// Apply pipe implementation
func (stgs pipe) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
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
		tmp, err := stg.Apply(ctx, r)
		// in case of any error
		// return immediately
		if err != nil {
			return r, err
		}
		// copy result back to
		// result structure
		r = tmp
	}
	return r, ctx.Err()
}
