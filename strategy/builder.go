package strategy

import (
	"fmt"

	"1pkg/gopium"
)

// list of registered types strategies
var (
	// comment annotation
	Note  gopium.StrategyName = "Comment_Note"
	Stamp gopium.StrategyName = "Comment_Stamp"
	// lexicographical sorts
	LexAsc  gopium.StrategyName = "Lexicographical_Ascending"
	LexDesc gopium.StrategyName = "Lexicographical_Descending"
	// pack for optimal mem util
	Pack gopium.StrategyName = "Memory_Pack"
	// explicit sys/type pads
	PadSys  gopium.StrategyName = "Explicit_Padings_System_Alignment"
	PadTnat gopium.StrategyName = "Explicit_Padings_Type_Natural"
	// false sharing guards
	FShareL1 gopium.StrategyName = "False_Sharing_CPU_L1"
	FShareL2 gopium.StrategyName = "False_Sharing_CPU_L2"
	FShareL3 gopium.StrategyName = "False_Sharing_CPU_L3"
	// cache line pad roundings
	CacheL1 gopium.StrategyName = "Cache_Rounding_CPU_L1"
	CacheL2 gopium.StrategyName = "Cache_Rounding_CPU_L2"
	CacheL3 gopium.StrategyName = "Cache_Rounding_CPU_L3"
	// start, end separate pads
	SepSys gopium.StrategyName = "Separate_Padding_System_Alignment"
	SepL1  gopium.StrategyName = "Separate_Padding_CPU_L1"
	SepL2  gopium.StrategyName = "Separate_Padding_CPU_L2"
	SepL3  gopium.StrategyName = "Separate_Padding_CPU_L3"
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
	// comment annotation
	case Note:
		return nt, nil
	case Stamp:
		return stmp, nil
	// lexicographical sorts
	case LexAsc:
		return lexasc, nil
	case LexAsc:
		return lexdesc, nil
	// pack for optimal mem util
	case Pack:
		return Pipe(
			filterpad,
			pck,
		), nil
	// explicit sys/type pads
	case PadSys:
		return Pipe(
			filterpad,
			padsys.C(b.c),
		), nil
	case PadTnat:
		return Pipe(
			filterpad,
			padtnat.C(b.c),
		), nil
	// false sharing guards
	case FShareL1:
		return Pipe(
			filterpad,
			fsharel1.C(b.c),
		), nil
	case FShareL2:
		return Pipe(
			filterpad,
			fsharel2.C(b.c),
		), nil
	case FShareL3:
		return Pipe(
			filterpad,
			fsharel3.C(b.c),
		), nil
	// cache line pad roundings
	case CacheL1:
		return cachel1.C(b.c), nil
	case CacheL2:
		return cachel2.C(b.c), nil
	case CacheL3:
		return cachel3.C(b.c), nil
	// start, end separate pads
	case SepSys:
		return sepsys.C(b.c), nil
	case SepL1:
		return sepl1.C(b.c), nil
	case SepL2:
		return sepl2.C(b.c), nil
	case SepL3:
		return sepl3.C(b.c), nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}

// Pipe concats list of strategy in one
// single piped strategy
func Pipe(stgs ...gopium.Strategy) gopium.Strategy {
	return pipe(stgs)
}
