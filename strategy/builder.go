package strategy

import (
	"fmt"
	"strings"

	"1pkg/gopium"
)

// list of registered types strategies
var (
	// comment annotation and others
	Nope  gopium.StrategyName = "nope"
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
	SepSysT gopium.StrategyName = "separate_padding_system_alignment_top"
	SepL1T  gopium.StrategyName = "separate_padding_cpu_l1_top"
	SepL2T  gopium.StrategyName = "separate_padding_cpu_l2_top"
	SepL3T  gopium.StrategyName = "separate_padding_cpu_l3_top"
	SepSysB gopium.StrategyName = "separate_padding_system_alignment_bottom"
	SepL1B  gopium.StrategyName = "separate_padding_cpu_l1_bottom"
	SepL2B  gopium.StrategyName = "separate_padding_cpu_l2_bottom"
	SepL3B  gopium.StrategyName = "separate_padding_cpu_l3_bottom"
	SepSysA gopium.StrategyName = "separate_padding_system_alignment_both"
	SepL1A  gopium.StrategyName = "separate_padding_cpu_l1_both"
	SepL2A  gopium.StrategyName = "separate_padding_cpu_l2_both"
	SepL3A  gopium.StrategyName = "separate_padding_cpu_l3_both"
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
	case Nope:
		return np, nil
	case Note:
		return nt, nil
	case Stamp:
		return stmp, nil
	case Group:
		return grp.Builder(b), nil
	// lexicographical and length sorts
	case LexAsc:
		return lexasc, nil
	case LexAsc:
		return lexdesc, nil
	case LenAsc:
		return lenasc, nil
	case LenAsc:
		return lendesc, nil
	// pack/unpack mem util
	case Pack:
		return Pipe(
			filterpad,
			pck,
		), nil
	case Unpack:
		return Pipe(
			filterpad,
			unpck,
		), nil
	// explicit sys/type pads
	case PadSys:
		return Pipe(
			filterpad,
			padsys.Curator(b.curator),
		), nil
	case PadTnat:
		return Pipe(
			filterpad,
			padtnat.Curator(b.curator),
		), nil
	// false sharing guards
	case FShareL1:
		return Pipe(
			filterpad,
			fsharel1.Curator(b.curator),
		), nil
	case FShareL2:
		return Pipe(
			filterpad,
			fsharel2.Curator(b.curator),
		), nil
	case FShareL3:
		return Pipe(
			filterpad,
			fsharel3.Curator(b.curator),
		), nil
	// cache line pad roundings
	case CacheL1:
		return cachel1.Curator(b.curator), nil
	case CacheL2:
		return cachel2.Curator(b.curator), nil
	case CacheL3:
		return cachel3.Curator(b.curator), nil
	// start, end separate pads
	case SepSysT:
		return sepsyst.Curator(b.curator), nil
	case SepL1T:
		return sepl1t.Curator(b.curator), nil
	case SepL2T:
		return sepl2t.Curator(b.curator), nil
	case SepL3T:
		return sepl3t.Curator(b.curator), nil
	case SepSysB:
		return sepsysb.Curator(b.curator), nil
	case SepL1B:
		return sepl1b.Curator(b.curator), nil
	case SepL2B:
		return sepl2b.Curator(b.curator), nil
	case SepL3B:
		return sepl3b.Curator(b.curator), nil
	case SepSysA:
		return Pipe(
			sepsyst.Curator(b.curator),
			sepsysb.Curator(b.curator),
		), nil
	case SepL1A:
		return Pipe(
			sepl1t.Curator(b.curator),
			sepl1b.Curator(b.curator),
		), nil
	case SepL2A:
		return Pipe(
			sepl2t.Curator(b.curator),
			sepl2b.Curator(b.curator),
		), nil
	case SepL3A:
		return Pipe(
			sepl3t.Curator(b.curator),
			sepl3b.Curator(b.curator),
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
func Tag(force bool, stgs ...gopium.StrategyName) gopium.Strategy {
	s := make([]string, 0, len(stgs))
	for _, stg := range stgs {
		s = append(s, string(stg))
	}
	return tag{tag: strings.Join(s, ","), force: force}
}
