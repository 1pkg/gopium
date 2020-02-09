package strategy

import (
	"fmt"

	"1pkg/gopium"
	"1pkg/gopium/types"
	gtypes "1pkg/gopium/types"
)

// List of registred types gopium.StrategyName
var (
	StrategyEnumerate gopium.StrategyName = "StrategyEnumerate"
)

// Builder defines types gopium.StrategyBuilder implementation
// that uses gtypes.Extractor as an extractor and related strategies
type Builder struct {
	extractor gtypes.Extractor
}

// NewBuilder creates instance of Builder
// and requires gtypes.Extractor to pass it to related strategies
func NewBuilder(extractor types.Extractor) Builder {
	return Builder{extractor: extractor}
}

// Build Builder implementation
func (b Builder) Build(name gopium.StrategyName) (gopium.Strategy, error) {
	switch name {
	case StrategyEnumerate:
		return stgenum(b), nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}
