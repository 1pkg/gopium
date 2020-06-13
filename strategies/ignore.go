package strategies

import (
	"context"

	"github.com/1pkg/gopium/gopium"
)

// list of ignore presets
var (
	ignr = ignore{}
)

// ignore defines nil strategy implementation
// that does nothing by returning original structure
type ignore struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Apply ignore implementation
func (stg ignore) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	return o, ctx.Err()
}
