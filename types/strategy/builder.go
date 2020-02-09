package strategy

import (
	"fmt"

	"1pkg/gopium"
	"1pkg/gopium/types"
	gtypes "1pkg/gopium/types"
)

// List of registred types gopium.StrategyName
var (
	Enumerate gopium.StrategyName = "Enumerate"
)

// Builder defines types gopium.StrategyBuilder implementation
// that uses gtypes.Extractor as an extractor and related strategies
type Builder struct {
	e gtypes.Extractor
}

// NewBuilder creates instance of Builder
// and requires gtypes.Extractor to pass it to related strategies
func NewBuilder(e types.Extractor) Builder {
	return Builder{e: e}
}

// Build Builder implementation
func (b Builder) Build(name gopium.StrategyName) (gopium.Strategy, error) {
	switch name {
	case Enumerate:
		return enumerate(b), nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}
