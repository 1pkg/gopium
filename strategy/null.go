package strategy

import (
	"context"

	"1pkg/gopium"
)

// list of null presets
var (
	nl = null{}
)

// null defines nil strategy implementation
// that only does nothing
type null struct{}

// Apply null implementation
func (stg null) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	return o, nil
}
