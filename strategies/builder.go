package strategies

import (
	"fmt"

	"1pkg/gopium"
)

// list of registered types strategies
const (
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
	// top, bottom separate pads
	SepSysT gopium.StrategyName = "separate_padding_system_alignment_top"
	SepL1T  gopium.StrategyName = "separate_padding_cpu_l1_top"
	SepL2T  gopium.StrategyName = "separate_padding_cpu_l2_top"
	SepL3T  gopium.StrategyName = "separate_padding_cpu_l3_top"
	SepSysB gopium.StrategyName = "separate_padding_system_alignment_bottom"
	SepL1B  gopium.StrategyName = "separate_padding_cpu_l1_bottom"
	SepL2B  gopium.StrategyName = "separate_padding_cpu_l2_bottom"
	SepL3B  gopium.StrategyName = "separate_padding_cpu_l3_bottom"
	// tag processors and modifiers
	PTGrp    gopium.StrategyName = "process_tag_group"
	AddTagS  gopium.StrategyName = "add_tag_group_soft"
	AddTagF  gopium.StrategyName = "add_tag_group_force"
	AddTagSD gopium.StrategyName = "add_tag_group_discrete"
	AddTagFD gopium.StrategyName = "add_tag_group_force_discrete"
	RmTagF   gopium.StrategyName = "remove_tag_group"
	// doc and comment annotations
	FNoteDoc  gopium.StrategyName = "doc_fields_annotate"
	FNoteCom  gopium.StrategyName = "comment_fields_annotate"
	StNoteDoc gopium.StrategyName = "doc_struct_annotate"
	StNoteCom gopium.StrategyName = "comment_struct_annotate"
	// lexicographical, length, embedded, exported sorts
	NLexAsc  gopium.StrategyName = "name_lexicographical_ascending"
	NLexDesc gopium.StrategyName = "name_lexicographical_descending"
	TLexAsc  gopium.StrategyName = "type_lexicographical_ascending"
	TLexDesc gopium.StrategyName = "type_lexicographical_descending"
	// filters and others
	FPad   gopium.StrategyName = "filter_pads"
	Ignore gopium.StrategyName = "ignore"
)

// Builder defines types gopium.StrategyBuilder implementation
// that uses gopium.Curator as an exposer and related strategies
type Builder struct {
	Curator gopium.Curator
}

// Build Builder implementation
func (b Builder) Build(names ...gopium.StrategyName) (gopium.Strategy, error) {
	// prepare result strategy pipe
	p := make(pipe, 0, len(names))
	for _, name := range names {
		var stg gopium.Strategy
		// build strategy by name
		switch name {
		// pack/unpack mem util
		case Pack:
			stg = pck
		case Unpack:
			stg = unpck
		// explicit sys/type pads
		case PadSys:
			stg = padsys.Curator(b.Curator)
		case PadTnat:
			stg = padtnat.Curator(b.Curator)
		// false sharing guards
		case FShareL1:
			stg = fsharel1.Curator(b.Curator)
		case FShareL2:
			stg = fsharel2.Curator(b.Curator)
		case FShareL3:
			stg = fsharel3.Curator(b.Curator)
		// cache line pad roundings
		case CacheL1:
			stg = cachel1.Curator(b.Curator)
		case CacheL2:
			stg = cachel2.Curator(b.Curator)
		case CacheL3:
			stg = cachel3.Curator(b.Curator)
		// top, bottom separate pads
		case SepSysT:
			stg = sepsyst.Curator(b.Curator)
		case SepL1T:
			stg = sepl1t.Curator(b.Curator)
		case SepL2T:
			stg = sepl2t.Curator(b.Curator)
		case SepL3T:
			stg = sepl3t.Curator(b.Curator)
		case SepSysB:
			stg = sepsysb.Curator(b.Curator)
		case SepL1B:
			stg = sepl1b.Curator(b.Curator)
		case SepL2B:
			stg = sepl2b.Curator(b.Curator)
		case SepL3B:
			stg = sepl3b.Curator(b.Curator)
		// tag processors and modifiers
		case PTGrp:
			stg = ptgrp.Builder(b)
		case AddTagS:
			stg = tags.Names(names...)
		case AddTagF:
			stg = tagf.Names(names...)
		case AddTagSD:
			stg = tagsd.Names(names...)
		case AddTagFD:
			stg = tagfd.Names(names...)
		case RmTagF:
			stg = tagf
		// doc and comment annotations
		case FNoteDoc:
			stg = fnotedoc
		case FNoteCom:
			stg = fnotecom
		case StNoteDoc:
			stg = stnotedoc
		case StNoteCom:
			stg = stnotecom
		// lexicographical, length, embedded, exported sorts
		case NLexAsc:
			stg = nlexasc
		case NLexDesc:
			stg = nlexdesc
		case TLexAsc:
			stg = tlexasc
		case TLexDesc:
			stg = tlexdesc
		// filters and others
		case FPad:
			stg = fpad
		case Ignore:
			stg = ignr
		default:
			return nil, fmt.Errorf("strategy %q wasn't found", name)
		}
		// append strategy to pipe
		p = append(p, stg)
	}
	return p, nil
}
