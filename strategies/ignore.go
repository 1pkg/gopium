package strategies

import (
	"context"

	"1pkg/gopium"
)

// list of ignore presets
var (
	ignr = ignore{}
)

// ignore defines nil strategy implementation
// that does nothing by returning original structure
type ignore struct{}

// Apply ignore implementation
func (stg ignore) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	return o, ctx.Err()
}
