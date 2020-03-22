package strategy

import (
	"fmt"
	"regexp"

	"1pkg/gopium"
)

// List of registered types gopium.StrategyName
var (
	Annotate        gopium.StrategyName = "Annotate"
	Stamp           gopium.StrategyName = "Stamp"
	FilterPad       gopium.StrategyName = "FilterPad"
	Lexicographical gopium.StrategyName = "Lexicographical"
	Memory          gopium.StrategyName = "Memory"
	PadType         gopium.StrategyName = "PadType"
	PadSys          gopium.StrategyName = "PadSys"
	CacheL1         gopium.StrategyName = "CacheL1"
	CacheL2         gopium.StrategyName = "CacheL2"
	CacheL3         gopium.StrategyName = "CacheL3"
	SeparateSys     gopium.StrategyName = "SeparateSys"
	SeparateL1      gopium.StrategyName = "SeparateL1"
	SeparateL2      gopium.StrategyName = "SeparateL2"
	SeparateL3      gopium.StrategyName = "SeparateL3"
	FalseShareL1    gopium.StrategyName = "FalseShareL1"
	FalseShareL2    gopium.StrategyName = "FalseShareL2"
	FalseShareL3    gopium.StrategyName = "FalseShareL3"
)

// Builder defines types gopium.StrategyBuilder implementation
// that uses gopium.Curator as an exposer and related strategies
type Builder struct {
	c gopium.Curator
}

// NewBuilder creates instance of Builder
// and requires gopium.Maven to pass it to related strategies
func NewBuilder(c gopium.Curator) Builder {
	return Builder{c: c}
}

// Build Builder implementation
func (b Builder) Build(name gopium.StrategyName) (gopium.Strategy, error) {
	// build strategy by name
	switch name {
	case Annotate:
		return annotate{}, nil
	case Stamp:
		return stamp{}, nil
	case FilterPad:
		regex, err := regexp.Compile(`^_$`)
		if err != nil {
			return nil, err
		}
		return filter{regex}, nil
	case Lexicographical:
		return lex{}, nil
	case Memory:
		return memory{}, nil
	case PadType:
		return pad{c: b.c, sys: false}, nil
	case PadSys:
		return pad{c: b.c, sys: true}, nil
	case CacheL1:
		return cache{c: b.c, l: 1}, nil
	case CacheL2:
		return cache{c: b.c, l: 2}, nil
	case CacheL3:
		return cache{c: b.c, l: 3}, nil
	case SeparateSys:
		return separate{c: b.c, sys: true}, nil
	case SeparateL1:
		return separate{c: b.c, l: 1}, nil
	case SeparateL2:
		return separate{c: b.c, l: 2}, nil
	case SeparateL3:
		return separate{c: b.c, l: 3}, nil
	case FalseShareL1:
		return fshare{c: b.c, l: 1}, nil
	case FalseShareL2:
		return fshare{c: b.c, l: 2}, nil
	case FalseShareL3:
		return fshare{c: b.c, l: 3}, nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}
