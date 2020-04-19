package strategies

import (
	"fmt"
	"strings"

	"1pkg/gopium"
)

// list of registered types strategies
const (
	// tag processors and modifiers
	TagGroup gopium.StrategyName = "process_tag_group"
	TagRm    gopium.StrategyName = "remove_tag_group"
	// filters and others
	FPad    gopium.StrategyName = "filter_pads"
	FEmb    gopium.StrategyName = "filter_embedded"
	FNotEmb gopium.StrategyName = "filter_not_embedded"
	FExp    gopium.StrategyName = "filter_exported"
	FNotExp gopium.StrategyName = "filter_not_exported"
	Nope    gopium.StrategyName = "nope"
	Void    gopium.StrategyName = "void"
	// lexicographical, length, embedded, exported sorts
	NLexAsc  gopium.StrategyName = "name_lexicographical_ascending"
	NLexDesc gopium.StrategyName = "name_lexicographical_descending"
	NLenAsc  gopium.StrategyName = "name_length_ascending"
	NLenDesc gopium.StrategyName = "name_length_descending"
	TLexAsc  gopium.StrategyName = "type_lexicographical_ascending"
	TLexDesc gopium.StrategyName = "type_lexicographical_descending"
	TLenAsc  gopium.StrategyName = "type_length_ascending"
	TLenDesc gopium.StrategyName = "type_length_descending"
	EmbAsc   gopium.StrategyName = "embedded_ascending"
	EmbDesc  gopium.StrategyName = "embedded_descending"
	ExpAsc   gopium.StrategyName = "exported_ascending"
	ExpDesc  gopium.StrategyName = "exported_descending"
	// pack/unpack mem util
	Pack   gopium.StrategyName = "memory_pack"
	Unpack gopium.StrategyName = "memory_unpack"
	// explicit sys/type pads
	PadSys  gopium.StrategyName = "explicit_padings_system_alignment"
	PadTnat gopium.StrategyName = "explicit_padings_type_natural"
	// false sharing guards
	FShareL1 gopium.StrategyName = "false_sharing_cpu_l1"
	FShareL2 gopium.StrategyName = "false_sharing_cpu_l2"
	FShareL3 gopium.StrategyName = "false_sharing_cpu_l3"
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
	// doc and comment annotations
	FNoteDoc  gopium.StrategyName = "doc_fields_annotate"
	FNoteCom  gopium.StrategyName = "comment_fields_annotate"
	StNoteDoc gopium.StrategyName = "doc_struct_annotate"
	StNoteCom gopium.StrategyName = "comment_struct_annotate"
	StampDoc  gopium.StrategyName = "doc_struct_stamp"
	StampCom  gopium.StrategyName = "comment_struct_stamp"
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
	// tag processors and modifiers
	case TagGroup:
		return grp.Builder(b), nil
	case TagRm:
		return tagrm, nil
	// filters and others
	case FPad:
		return fpad, nil
	case FEmb:
		return femb, nil
	case FNotEmb:
		return fnotemb, nil
	case FExp:
		return fexp, nil
	case FNotExp:
		return fnotexp, nil
	case Nope:
		return np, nil
	case Void:
		return vd, nil
	// lexicographical, length, embedded, exported sorts
	case NLexAsc:
		return nlexasc, nil
	case NLexDesc:
		return nlexdesc, nil
	case NLenAsc:
		return nlenasc, nil
	case NLenDesc:
		return nlendesc, nil
	case TLexAsc:
		return tlexasc, nil
	case TLexDesc:
		return tlexdesc, nil
	case TLenAsc:
		return tlenasc, nil
	case TLenDesc:
		return tlendesc, nil
	case EmbAsc:
		return embasc, nil
	case EmbDesc:
		return embdesc, nil
	case ExpAsc:
		return expasc, nil
	case ExpDesc:
		return expdesc, nil
	// pack/unpack mem util
	case Pack:
		return pck, nil
	case Unpack:
		return unpck, nil
	// explicit sys/type pads
	case PadSys:
		return padsys.Curator(b.curator), nil
	case PadTnat:
		return padtnat.Curator(b.curator), nil
	// false sharing guards
	case FShareL1:
		return fsharel1.Curator(b.curator), nil
	case FShareL2:
		return fsharel2.Curator(b.curator), nil
	case FShareL3:
		return fsharel3.Curator(b.curator), nil
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
	// doc and comment annotations
	case FNoteDoc:
		return fnotedoc, nil
	case FNoteCom:
		return fnotecom, nil
	case StNoteDoc:
		return stnotedoc, nil
	case StNoteCom:
		return stnotecom, nil
	case StampDoc:
		return stampdoc, nil
	case StampCom:
		return stampcom, nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}

// Pipe helps to concat slice of strategies
// in one single pipe strategy
func Pipe(stgs ...gopium.Strategy) gopium.Strategy {
	return pipe(stgs)
}

// Tag concats slice of strategy names
// in one single tag strategy
func Tag(group string, force, discrete bool, stgs ...gopium.StrategyName) gopium.Strategy {
	s := make([]string, 0, len(stgs))
	for _, stg := range stgs {
		s = append(s, string(stg))
	}
	return tag{
		tag:      strings.Join(s, ","),
		group:    group,
		force:    force,
		discrete: discrete,
	}
}
