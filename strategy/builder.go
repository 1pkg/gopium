package strategy

import (
	"fmt"
	"regexp"

	"1pkg/gopium"
)

// List of registered types gopium.StrategyName
var (
	Enumerate              gopium.StrategyName = "Enumerate"
	FilterPad              gopium.StrategyName = "FilterPad"
	Lexicographical        gopium.StrategyName = "Lexicographical"
	Memory                 gopium.StrategyName = "Memory"
	PaddingType            gopium.StrategyName = "PaddingType"
	PaddingSys             gopium.StrategyName = "PaddingSys"
	CachingCPUL1           gopium.StrategyName = "CachingCPUL1"
	CachingCPUL2           gopium.StrategyName = "CachingCPUL2"
	CachingCPUL3           gopium.StrategyName = "CachingCPUL3"
	SeparatingSys          gopium.StrategyName = "SeparatingSys"
	SeparatingCachingCPUL1 gopium.StrategyName = "SeparatingCachingCPUL1"
	SeparatingCachingCPUL2 gopium.StrategyName = "SeparatingCachingCPUL2"
	SeparatingCachingCPUL3 gopium.StrategyName = "SeparatingCachingCPUL3"
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
	case PaddingType:
		stg = padding{m: b.m, sys: false}
	case PaddingSys:
		stg = padding{m: b.m, sys: true}
	case PaddingSys:
		stg = padding{m: b.m, sys: true}
	case CachingCPUL1:
		stg = caching{m: b.m, l: 1}
	case CachingCPUL2:
		stg = caching{m: b.m, l: 2}
	case CachingCPUL3:
		stg = caching{m: b.m, l: 3}
	case SeparatingSys:
		stg = separating{b.m}
	case SeparatingCachingCPUL1:
		stg = separating_caching{m: b.m, l: 1}
	case SeparatingCachingCPUL2:
		stg = separating_caching{m: b.m, l: 2}
	case SeparatingCachingCPUL2:
		stg = separating_caching{m: b.m, l: 3}
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
