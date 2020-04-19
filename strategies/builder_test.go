package strategies

import (
	"fmt"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/mocks"
)

func TestBuilder(t *testing.T) {
	// prepare
	b := NewBuilder(mocks.Maven{})
	table := map[string]struct {
		name gopium.StrategyName
		stg  gopium.Strategy
		err  error
	}{
		"empty name should lead to builder error": {
			err: fmt.Errorf(`strategy "" wasn't found`),
		},
		"unknown name should lead to builder error": {
			name: "test",
			err:  fmt.Errorf(`strategy "test" wasn't found`),
		},
		"`process_tag_group` name should lead to relevant strategy": {
			name: TagGroup,
			stg:  grp.Builder(b),
		},
		"`remove_tag_group` name should lead to relevant strategy": {
			name: TagRm,
			stg:  tagrm,
		},
		"`filter_pads` name should lead to relevant strategy": {
			name: FPad,
			stg:  fpad,
		},
		"`filter_embedded` name should lead to relevant strategy": {
			name: FEmb,
			stg:  femb,
		},
		"`filter_not_embedded` name should lead to relevant strategy": {
			name: FNotEmb,
			stg:  fnotemb,
		},
		"`filter_exported` name should lead to relevant strategy": {
			name: FExp,
			stg:  fexp,
		},
		"`filter_not_exported` name should lead to relevant strategy": {
			name: FNotExp,
			stg:  fnotexp,
		},
		"`nope` name should lead to relevant strategy": {
			name: Nope,
			stg:  np,
		},
		"`void` name should lead to relevant strategy": {
			name: Void,
			stg:  vd,
		},
		"`name_lexicographical_ascending` name should lead to relevant strategy": {
			name: NLexAsc,
			stg:  nlexasc,
		},
		"`name_lexicographical_descending` name should lead to relevant strategy": {
			name: NLexDesc,
			stg:  nlexdesc,
		},
		"`name_length_ascending` name should lead to relevant strategy": {
			name: NLenAsc,
			stg:  nlenasc,
		},
		"`name_length_descending` name should lead to relevant strategy": {
			name: NLenDesc,
			stg:  nlendesc,
		},
		"`type_lexicographical_ascending` name should lead to relevant strategy": {
			name: TLexAsc,
			stg:  tlexasc,
		},
		"`type_lexicographical_descending` name should lead to relevant strategy": {
			name: TLexDesc,
			stg:  tlexdesc,
		},
		"`type_length_ascending` name should lead to relevant strategy": {
			name: TLenAsc,
			stg:  tlenasc,
		},
		"`type_length_descending` name should lead to relevant strategy": {
			name: TLenDesc,
			stg:  tlendesc,
		},
		"`embedded_ascending` name should lead to relevant strategy": {
			name: EmbAsc,
			stg:  embasc,
		},
		"`embedded_descending` name should lead to relevant strategy": {
			name: EmbDesc,
			stg:  embdesc,
		},
		"`exported_ascending` name should lead to relevant strategy": {
			name: ExpAsc,
			stg:  expasc,
		},
		"`exported_descending` name should lead to relevant strategy": {
			name: ExpDesc,
			stg:  expdesc,
		},
		"`memory_pack` name should lead to relevant strategy": {
			name: Pack,
			stg:  pck,
		},
		"`memory_unpack` name should lead to relevant strategy": {
			name: Unpack,
			stg:  unpck,
		},
		"`explicit_padings_system_alignment` name should lead to relevant strategy": {
			name: PadSys,
			stg:  padsys.Curator(b.curator),
		},
		"`explicit_padings_type_natural` name should lead to relevant strategy": {
			name: PadTnat,
			stg:  padtnat.Curator(b.curator),
		},
		"`false_sharing_cpu_l1` name should lead to relevant strategy": {
			name: FShareL1,
			stg:  fsharel1.Curator(b.curator),
		},
		"`false_sharing_cpu_l2` name should lead to relevant strategy": {
			name: FShareL2,
			stg:  fsharel2.Curator(b.curator),
		},
		"`false_sharing_cpu_l3` name should lead to relevant strategy": {
			name: FShareL3,
			stg:  fsharel3.Curator(b.curator),
		},
		"`cache_rounding_cpu_l1` name should lead to relevant strategy": {
			name: CacheL1,
			stg:  cachel1.Curator(b.curator),
		},
		"`cache_rounding_cpu_l2` name should lead to relevant strategy": {
			name: CacheL2,
			stg:  cachel2.Curator(b.curator),
		},
		"`cache_rounding_cpu_l3` name should lead to relevant strategy": {
			name: CacheL3,
			stg:  cachel3.Curator(b.curator),
		},
		"`separate_padding_system_alignment_top` name should lead to relevant strategy": {
			name: SepSysT,
			stg:  sepsyst.Curator(b.curator),
		},
		"`separate_padding_cpu_l1_top` name should lead to relevant strategy": {
			name: SepL1T,
			stg:  sepl1t.Curator(b.curator),
		},
		"`separate_padding_cpu_l2_top` name should lead to relevant strategy": {
			name: SepL2T,
			stg:  sepl2t.Curator(b.curator),
		},
		"`separate_padding_cpu_l3_top` name should lead to relevant strategy": {
			name: SepL3T,
			stg:  sepl3t.Curator(b.curator),
		},
		"`separate_padding_system_alignment_bottom` name should lead to relevant strategy": {
			name: SepSysB,
			stg:  sepsysb.Curator(b.curator),
		},
		"`separate_padding_cpu_l1_bottom` name should lead to relevant strategy": {
			name: SepL1B,
			stg:  sepl1b.Curator(b.curator),
		},
		"`separate_padding_cpu_l2_bottom` name should lead to relevant strategy": {
			name: SepL2B,
			stg:  sepl2b.Curator(b.curator),
		},
		"`separate_padding_cpu_l3_bottom` name should lead to relevant strategy": {
			name: SepL3B,
			stg:  sepl3b.Curator(b.curator),
		},
		"`doc_fields_annotate` name should lead to relevant strategy": {
			name: FNoteDoc,
			stg:  fnotedoc,
		},
		"`comment_fields_annotate` name should lead to relevant strategy": {
			name: FNoteCom,
			stg:  fnotecom,
		},
		"`doc_struct_annotate` name should lead to relevant strategy": {
			name: StNoteDoc,
			stg:  stnotedoc,
		},
		"`comment_struct_annotate` name should lead to relevant strategy": {
			name: StNoteCom,
			stg:  stnotecom,
		},
		"`doc_struct_stamp` name should lead to relevant strategy": {
			name: StampDoc,
			stg:  stampdoc,
		},
		"`comment_struct_stamp` name should lead to relevant strategy": {
			name: StampCom,
			stg:  stampcom,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			stg, err := b.Build(tcase.name)
			// check
			if !reflect.DeepEqual(stg, tcase.stg) {
				t.Errorf("actual %v doesn't equal to expected %v", stg, tcase.stg)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}

func TestPipeBuilder(t *testing.T) {
	// prepare
	table := map[string]struct {
		stgs []gopium.Strategy
		stg  gopium.Strategy
	}{
		"nil strategies should lead to nil pipe": {
			stg: pipe(nil),
		},
		"empty strategies should lead to empty pipe": {
			stgs: []gopium.Strategy{},
			stg:  pipe{},
		},
		"non empty strategies should lead actual pipe": {
			stgs: []gopium.Strategy{pck, unpck, vd},
			stg:  pipe([]gopium.Strategy{pck, unpck, vd}),
		},
		"nested strategies should lead to actual pipe": {
			stgs: []gopium.Strategy{pck, pipe([]gopium.Strategy{pck, unpck, vd}), vd},
			stg:  pipe([]gopium.Strategy{pck, pipe([]gopium.Strategy{pck, unpck, vd}), vd}),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			stg := Pipe(tcase.stgs...)
			// check
			if !reflect.DeepEqual(stg, tcase.stg) {
				t.Errorf("actual %+v doesn't equal to expected %+v", stg, tcase.stg)
			}
		})
	}
}

func TestTagBuilder(t *testing.T) {
	// prepare
	table := map[string]struct {
		group    string
		force    bool
		discrete bool
		stgs     []gopium.StrategyName
		stg      gopium.Strategy
	}{
		"nil strategies should lead to nil tag": {
			group:    "test",
			force:    true,
			discrete: true,
			stg: tag{
				tag:      "",
				group:    "test",
				force:    true,
				discrete: true,
			},
		},
		"empty strategies should lead to empty tag": {
			group: "test",
			force: true,
			stgs:  []gopium.StrategyName{},
			stg: tag{
				tag:   "",
				group: "test",
				force: true,
			},
		},
		"non empty strategies should lead actual tag": {
			group: "test",
			stgs:  []gopium.StrategyName{Pack, Unpack, Void},
			stg: tag{
				tag:   "memory_pack,memory_unpack,void",
				group: "test",
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			stg := Tag(tcase.group, tcase.force, tcase.discrete, tcase.stgs...)
			// check
			if !reflect.DeepEqual(stg, tcase.stg) {
				t.Errorf("actual %+v doesn't equal to expected %+v", stg, tcase.stg)
			}
		})
	}
}
