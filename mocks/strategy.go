package mocks

import (
	"context"

	"1pkg/gopium"
)

// StrategyMock defines mock strategy implementation
type StrategyMock struct {
	R   gopium.Struct
	Err error
}

// Apply mock implementation
func (stg StrategyMock) Apply(context.Context, gopium.Struct) (gopium.Struct, error) {
	return stg.R, stg.Err
}
