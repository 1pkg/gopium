package strategies

import (
	"context"

	"1pkg/gopium"
)

// list of discard presets
var (
	dis = discard{}
)

// nope defines discard strategy implementation
// that discards struct fields by returning void struct
type discard struct{}

// Apply discard implementation
func (stg discard) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	return gopium.Struct{}, ctx.Err()
}
