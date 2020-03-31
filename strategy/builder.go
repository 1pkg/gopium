package strategy

import (
	"fmt"
	"strings"

	"1pkg/gopium"
)

// list of registered types strategies
var (
	// comment annotation and others
	Nil   gopium.StrategyName = "nil"
	Note  gopium.StrategyName = "comment_fields_annotate"
	Stamp gopium.StrategyName = "comment_struct_stamp"
	Group gopium.StrategyName = "group_tag"
	// lexicographical and length sorts
	LexAsc  gopium.StrategyName = "lexicographical_ascending"
	LexDesc gopium.StrategyName = "lexicographical_descending"
	LenAsc  gopium.StrategyName = "length_ascending"
	LenDesc gopium.StrategyName = "length_descending"
	// pack/unpack mem util
	Pack   gopium.StrategyName = "memory_pack"
	Unpack gopium.StrategyName = "memory_unpack"
	// explicit sys/type pads
	PadSys  gopium.StrategyName = "explicit_padings_system_alignment"
	PadTnat gopium.StrategyName = "explicit_padings_type_natural"
	// false sharing guards
	FShareL1 gopium.StrategyName = "false_sharing_cpu_l1"
	FShareL2 gopium.StrategyName = "false_sharing_cpu_l2"
	FShareL3 gopium.StrategyName = "false_sharing_cpu_l2"
	// cache line pad roundings
	CacheL1 gopium.StrategyName = "cache_rounding_cpu_l1"
	CacheL2 gopium.StrategyName = "cache_rounding_cpu_l2"
	CacheL3 gopium.StrategyName = "cache_rounding_cpu_l3"
	// start, end separate pads
	SepSys gopium.StrategyName = "separate_padding_system_alignment"
	SepL1  gopium.StrategyName = "separate_padding_cpu_l1"
	SepL2  gopium.StrategyName = "separate_padding_cpu_l2"
	SepL3  gopium.StrategyName = "separate_padding_cpu_l3"
)

// Builder defines types gopium.StrategyBuilder implementation
// that uses gopium.Curator as an exposer and related strategies
type Builder struct {
	curator gopium.Curator
}

// NewBuilder creates instance of Builder
// and requires gopium.Maven to pass it to related strategies
func NewBuilder(curator gopium.Curator) Builder {
	return Builder{curator: curator}
}

// Build Builder implementation
func (b Builder) Build(name gopium.StrategyName) (gopium.Strategy, error) {
	// build strategy by name
	switch name {
	// comment annotation and others
	case Nil:
		return nl, nil
	case Note:
		return nt, nil
	case Stamp:
		return stmp, nil
	case Group:
		return grp.Builder(b), nil
	// lexicographical and length sorts
	case LexAsc:
		return Pipe(
			lexasc,
			taglexasc,
		), nil
	case LexAsc:
		return Pipe(
			lexdesc,
			taglexdesc,
		), nil
	case LenAsc:
		return Pipe(
			lenasc,
			taglenasc,
		), nil
	case LenAsc:
		return Pipe(
			lendesc,
			taglendesc,
		), nil
	// pack/unpack mem util
	case Pack:
		return Pipe(
			filterpad,
			pck,
			tagpack,
		), nil
	case Unpack:
		return Pipe(
			filterpad,
			unpck,
			tagunpack,
		), nil
	// explicit sys/type pads
	case PadSys:
		return Pipe(
			filterpad,
			padsys.Curator(b.curator),
			tagpadsys,
		), nil
	case PadTnat:
		return Pipe(
			filterpad,
			padtnat.Curator(b.curator),
			tagpadtnat,
		), nil
	// false sharing guards
	case FShareL1:
		return Pipe(
			filterpad,
			fsharel1.Curator(b.curator),
			tagfsahrel1,
		), nil
	case FShareL2:
		return Pipe(
			filterpad,
			fsharel2.Curator(b.curator),
			tagfsahrel2,
		), nil
	case FShareL3:
		return Pipe(
			filterpad,
			fsharel3.Curator(b.curator),
			tagfsahrel3,
		), nil
	// cache line pad roundings
	case CacheL1:
		return Pipe(
			cachel1.Curator(b.curator),
			tagcachel1,
		), nil
	case CacheL2:
		return Pipe(
			cachel2.Curator(b.curator),
			tagcachel2,
		), nil
	case CacheL3:
		return Pipe(
			cachel3.Curator(b.curator),
			tagcachel3,
		), nil
	// start, end separate pads
	case SepSys:
		return Pipe(
			sepsys.Curator(b.curator),
			tagsepsys,
		), nil
	case SepL1:
		return Pipe(
			sepl1.Curator(b.curator),
			tagsepl1,
		), nil
	case SepL2:
		return Pipe(
			sepl2.Curator(b.curator),
			tagsepl2,
		), nil
	case SepL3:
		return Pipe(
			sepl3.Curator(b.curator),
			tagsepl3,
		), nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}

// Pipe helps to concat list of strategies
// in one single pipe strategy
func Pipe(stgs ...gopium.Strategy) gopium.Strategy {
	return pipe(stgs)
}

// Tag concats list of strategy names
// in one single tag strategy
func Tag(stgs ...gopium.StrategyName) gopium.Strategy {
	s := make([]string, 0, len(stgs))
	for _, stg := range stgs {
		s = append(s, string(stg))
	}
	return tag{tag: strings.Join(s, ","), force: true}
}
