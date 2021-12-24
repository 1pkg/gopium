package strategies

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/1pkg/gopium/gopium"
)

// list of registered strategies names
const (
	// pack/unpack mem util
	Pack   gopium.StrategyName = "memory_pack"
	Unpack gopium.StrategyName = "memory_unpack"
	// explicit sys/type pads
	PadSys  gopium.StrategyName = "explicit_paddings_system_alignment"
	PadTnat gopium.StrategyName = "explicit_paddings_type_natural"
	// false sharing guards
	FShareL1 gopium.StrategyName = "false_sharing_cpu_l1"
	FShareL2 gopium.StrategyName = "false_sharing_cpu_l2"
	FShareL3 gopium.StrategyName = "false_sharing_cpu_l3"
	FShareB  gopium.StrategyName = "false_sharing_bytes_%d"
	// cache line pad roundings
	CacheL1D gopium.StrategyName = "cache_rounding_cpu_l1_discrete"
	CacheL2D gopium.StrategyName = "cache_rounding_cpu_l2_discrete"
	CacheL3D gopium.StrategyName = "cache_rounding_cpu_l3_discrete"
	CacheBD  gopium.StrategyName = "cache_rounding_bytes_%d_discrete"
	CacheL1F gopium.StrategyName = "cache_rounding_cpu_l1_full"
	CacheL2F gopium.StrategyName = "cache_rounding_cpu_l2_full"
	CacheL3F gopium.StrategyName = "cache_rounding_cpu_l3_full"
	CacheBF  gopium.StrategyName = "cache_rounding_bytes_%d_full"
	// top, bottom separate pads
	SepSysT gopium.StrategyName = "separate_padding_system_alignment_top"
	SepSysB gopium.StrategyName = "separate_padding_system_alignment_bottom"
	SepL1T  gopium.StrategyName = "separate_padding_cpu_l1_top"
	SepL2T  gopium.StrategyName = "separate_padding_cpu_l2_top"
	SepL3T  gopium.StrategyName = "separate_padding_cpu_l3_top"
	SepBT   gopium.StrategyName = "separate_padding_bytes_%d_top"
	SepL1B  gopium.StrategyName = "separate_padding_cpu_l1_bottom"
	SepL2B  gopium.StrategyName = "separate_padding_cpu_l2_bottom"
	SepL3B  gopium.StrategyName = "separate_padding_cpu_l3_bottom"
	SepBB   gopium.StrategyName = "separate_padding_bytes_%d_bottom"
	// tag processors and modifiers
	ProcTag  gopium.StrategyName = "process_tag_group"
	AddTagS  gopium.StrategyName = "add_tag_group_soft"
	AddTagF  gopium.StrategyName = "add_tag_group_force"
	AddTagSD gopium.StrategyName = "add_tag_group_discrete"
	AddTagFD gopium.StrategyName = "add_tag_group_force_discrete"
	RmTagF   gopium.StrategyName = "remove_tag_group"
	// doc and comment annotations
	FNoteDoc  gopium.StrategyName = "fields_annotate_doc"
	FNoteCom  gopium.StrategyName = "fields_annotate_comment"
	StNoteDoc gopium.StrategyName = "struct_annotate_doc"
	StNoteCom gopium.StrategyName = "struct_annotate_comment"
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
	Curator gopium.Curator `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; struct ptr scan size: 16 bytes; - 🌺 gopium @1pkg

// Build Builder implementation
func (b Builder) Build(names ...gopium.StrategyName) (gopium.Strategy, error) {
	// prepare result strategy pipe
	p := make(pipe, 0, len(names))
	for _, name := range names {
		var stg gopium.Strategy
		// build strategy by name
		switch {
		// pack/unpack mem util
		case b.marchp(name, Pack):
			stg = pck
		case b.marchp(name, Unpack):
			stg = unpck
		// explicit sys/type pads
		case b.marchp(name, PadSys):
			stg = padsys.Curator(b.Curator)
		case b.marchp(name, PadTnat):
			stg = padtnat.Curator(b.Curator)
		// false sharing guards
		case b.marchp(name, FShareL1):
			stg = fsharel1.Curator(b.Curator)
		case b.marchp(name, FShareL2):
			stg = fsharel2.Curator(b.Curator)
		case b.marchp(name, FShareL3):
			stg = fsharel3.Curator(b.Curator)
		case b.marchp(name, FShareB):
			var bytes uint
			if err := b.scanp(name, FShareB, &bytes); err != nil {
				return nil, err
			}
			stg = fshareb.Bytes(bytes).Curator(b.Curator)
		// cache line pad roundings
		case b.marchp(name, CacheL1D):
			stg = cachel1d.Curator(b.Curator)
		case b.marchp(name, CacheL2D):
			stg = cachel2d.Curator(b.Curator)
		case b.marchp(name, CacheL3D):
			stg = cachel3d.Curator(b.Curator)
		case b.marchp(name, CacheBD):
			var bytes uint
			if err := b.scanp(name, CacheBD, &bytes); err != nil {
				return nil, err
			}
			stg = cachebd.Bytes(bytes).Curator(b.Curator)
		case b.marchp(name, CacheL1F):
			stg = cachel1f.Curator(b.Curator)
		case b.marchp(name, CacheL2F):
			stg = cachel2f.Curator(b.Curator)
		case b.marchp(name, CacheL3F):
			stg = cachel3f.Curator(b.Curator)
		case b.marchp(name, CacheBF):
			var bytes uint
			if err := b.scanp(name, CacheBF, &bytes); err != nil {
				return nil, err
			}
			stg = cachebf.Bytes(bytes).Curator(b.Curator)
		// top, bottom separate pads
		case b.marchp(name, SepSysT):
			stg = sepsyst.Curator(b.Curator)
		case b.marchp(name, SepSysB):
			stg = sepsysb.Curator(b.Curator)
		case b.marchp(name, SepL1T):
			stg = sepl1t.Curator(b.Curator)
		case b.marchp(name, SepL2T):
			stg = sepl2t.Curator(b.Curator)
		case b.marchp(name, SepL3T):
			stg = sepl3t.Curator(b.Curator)
		case b.marchp(name, SepBT):
			var bytes uint
			if err := b.scanp(name, SepBT, &bytes); err != nil {
				return nil, err
			}
			stg = sepbt.Bytes(bytes).Curator(b.Curator)
		case b.marchp(name, SepL1B):
			stg = sepl1b.Curator(b.Curator)
		case b.marchp(name, SepL2B):
			stg = sepl2b.Curator(b.Curator)
		case b.marchp(name, SepL3B):
			stg = sepl3b.Curator(b.Curator)
		case b.marchp(name, SepBB):
			var bytes uint
			if err := b.scanp(name, SepBB, &bytes); err != nil {
				return nil, err
			}
			stg = sepbb.Bytes(bytes).Curator(b.Curator)
		// tag processors and modifiers
		case b.marchp(name, ProcTag):
			stg = ptag.Builder(b)
		case b.marchp(name, AddTagS):
			stg = tags.Names(names...)
		case b.marchp(name, AddTagF):
			stg = tagf.Names(names...)
		case b.marchp(name, AddTagSD):
			stg = tagsd.Names(names...)
		case b.marchp(name, AddTagFD):
			stg = tagfd.Names(names...)
		case b.marchp(name, RmTagF):
			stg = tagf
		// doc and comment annotations
		case b.marchp(name, FNoteDoc):
			stg = fnotedoc
		case b.marchp(name, FNoteCom):
			stg = fnotecom
		case b.marchp(name, StNoteDoc):
			stg = stnotedoc
		case b.marchp(name, StNoteCom):
			stg = stnotecom
		// lexicographical, length, embedded, exported sorts
		case b.marchp(name, NLexAsc):
			stg = nlexasc
		case b.marchp(name, NLexDesc):
			stg = nlexdesc
		case b.marchp(name, TLexAsc):
			stg = tlexasc
		case b.marchp(name, TLexDesc):
			stg = tlexdesc
		// filters and others
		case b.marchp(name, FPad):
			stg = fpad
		case b.marchp(name, Ignore):
			stg = ignr
		default:
			return nil, fmt.Errorf("strategy %q wasn't found", name)
		}
		// append strategy to pipe
		p = append(p, stg)
	}
	return p, nil
}

// marchp checks if strahtegy name matches pattern
func (b Builder) marchp(name gopium.StrategyName, pattern gopium.StrategyName) bool {
	// for matching we need to use regex
	// note only: %d %f %s formatters are used
	p := string(pattern)
	p = strings.ReplaceAll(p, "%d", ".*?")
	p = strings.ReplaceAll(p, "%f", ".*?")
	p = strings.ReplaceAll(p, "%s", ".*?")
	p = fmt.Sprintf("^%s$", p)
	// try to compile artificial regex
	regex, err := regexp.Compile(p)
	return err == nil && regex.MatchString(string(name))
}

// scanp scans name with provided pattern to variable list
func (b Builder) scanp(name gopium.StrategyName, pattern gopium.StrategyName, vars ...interface{}) error {
	// prepare strings to be scanned
	p := strings.ReplaceAll(string(pattern), "_", " ")
	n := strings.ReplaceAll(string(name), "_", " ")
	// perform the scan and handle errors
	if _, err := fmt.Sscanf(n, p, vars...); err != nil {
		return fmt.Errorf("pattern %q can't be scanned for strategy %q %v", pattern, name, err)
	}
	return nil
}
