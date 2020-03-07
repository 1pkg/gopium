package strategy

import (
	"fmt"

	"1pkg/gopium"
	"1pkg/gopium/types"
	gtypes "1pkg/gopium/types"
)

// List of registered types gopium.StrategyName
var (
	Enumerate       gopium.StrategyName = "Enumerate"
	Lexicographical gopium.StrategyName = "Lexicographical"
	Memory          gopium.StrategyName = "Memory"
)

// List of registered modes gopium.StrategyMode
const (
	WithNone gopium.StrategyMode = 1 << iota
	WithAnnotation
	WithStamp
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
func (b Builder) Build(name gopium.StrategyName, mode gopium.StrategyMode) (gopium.Strategy, error) {
	// build strategy by name
	var stg gopium.Strategy
	switch name {
	case Enumerate:
		stg = enum{b.extractor}
	case Lexicographical:
		stg = lexicographical{b.extractor}
	case Memory:
		stg = memory{b.extractor}
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
	// iterate through all registered modes
	for mask := WithNone; mask != WithStamp; mask = mask << 1 {
		// in case mode doesn't have current mask
		// just skip current mask
		if !mode.Has(mask) {
			continue
		}
		// otherwise apply registered mode
		switch mask {
		case WithAnnotation:
			stg = annotate{stg}
		case WithStamp:
			stg = stamp{stg}
		}
	}
	return stg, nil
}
