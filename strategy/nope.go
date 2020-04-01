package strategy

import (
	"context"

	"1pkg/gopium"
)

// list of nope presets
var (
	np = nope{}
)

// nope defines nil strategy implementation
// that only does nothing and just returns original
type nope struct{}

// Apply nope implementation
func (stg nope) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	return o, nil
}
