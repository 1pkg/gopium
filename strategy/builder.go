package strategy

import (
	"fmt"
	"regexp"

	"1pkg/gopium"
)

// List of registered types gopium.StrategyName
var (
	Enumerate       gopium.StrategyName = "Enumerate"
	FilterPad       gopium.StrategyName = "FilterPad"
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
// that uses gopium.Maven as an exposer and related strategies
type Builder struct {
	m gopium.Maven
}

// NewBuilder creates instance of Builder
// and requires gopium.Maven to pass it to related strategies
func NewBuilder(m gopium.Maven) Builder {
	return Builder{m: m}
}

// Build Builder implementation
func (b Builder) Build(name gopium.StrategyName, mode gopium.StrategyMode) (gopium.Strategy, error) {
	// build strategy by name
	var stg gopium.Strategy
	switch name {
	case Enumerate:
		stg = enum{b.m}
	case FilterPad:
		regex, err := regexp.Compile(`^_$`)
		if err != nil {
			return nil, err
		}
		stg = filter{m: b.m, r: regex}
	case Lexicographical:
		stg = lexicographical{b.m}
	case Memory:
		stg = memory{b.m}
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
