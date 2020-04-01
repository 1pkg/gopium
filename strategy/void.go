package strategy

import (
	"context"

	"1pkg/gopium"
)

// list of void presets
var (
	vd = void{}
)

// nope defines void strategy implementation
// that only does nothing and just returns void struct
type void struct{}

// Apply void implementation
func (stg void) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	return gopium.Struct{}, nil
}
