package mocks

import (
	"context"

	"1pkg/gopium"
)

// Strategy defines mock strategy implementation
type Strategy struct {
	R   gopium.Struct
	Err error
}

// Apply mock implementation
func (stg Strategy) Apply(context.Context, gopium.Struct) (gopium.Struct, error) {
	return stg.R, stg.Err
}