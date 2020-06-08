package strategies

import (
	"fmt"
	"reflect"
	"testing"

	"1pkg/gopium/gopium"
	"1pkg/gopium/tests/mocks"
)

func TestBuilder(t *testing.T) {
	// prepare
	b := Builder{Curator: mocks.Maven{}}
	table := map[string]struct {
		names []gopium.StrategyName
		stg   gopium.Strategy
		err   error
	}{
		// pack/unpack mem util
		"`memory_pack` name should return expected strategy": {
			names: []gopium.StrategyName{Pack},
			stg:   pipe([]gopium.Strategy{pck}),
		},
		"`memory_unpack` name should return expected strategy": {
			names: []gopium.StrategyName{Unpack},
			stg:   pipe([]gopium.Strategy{unpck}),
		},
		// explicit sys/type pads
		"`explicit_padings_system_alignment` name should return expected strategy": {
			names: []gopium.StrategyName{PadSys},
			stg:   pipe([]gopium.Strategy{padsys.Curator(b.Curator)}),
		},
		"`explicit_padings_type_natural` name should return expected strategy": {
			names: []gopium.StrategyName{PadTnat},
			stg:   pipe([]gopium.Strategy{padtnat.Curator(b.Curator)}),
		},
		// false sharing guards
		"`false_sharing_cpu_l1` name should return expected strategy": {
			names: []gopium.StrategyName{FShareL1},
			stg:   pipe([]gopium.Strategy{fsharel1.Curator(b.Curator)}),
		},
		"`false_sharing_cpu_l2` name should return expected strategy": {
			names: []gopium.StrategyName{FShareL2},
			stg:   pipe([]gopium.Strategy{fsharel2.Curator(b.Curator)}),
		},
		"`false_sharing_cpu_l3` name should return expected strategy": {
			names: []gopium.StrategyName{FShareL3},
			stg:   pipe([]gopium.Strategy{fsharel3.Curator(b.Curator)}),
		},
		// cache line pad roundings
		"`cache_rounding_cpu_l1` name should return expected strategy": {
			names: []gopium.StrategyName{CacheL1},
			stg:   pipe([]gopium.Strategy{cachel1.Curator(b.Curator)}),
		},
		"`cache_rounding_cpu_l2` name should return expected strategy": {
			names: []gopium.StrategyName{CacheL2},
			stg:   pipe([]gopium.Strategy{cachel2.Curator(b.Curator)}),
		},
		"`cache_rounding_cpu_l3` name should return expected strategy": {
			names: []gopium.StrategyName{CacheL3},
			stg:   pipe([]gopium.Strategy{cachel3.Curator(b.Curator)}),
		},
		"`full_cache_rounding_cpu_l1` name should return expected strategy": {
			names: []gopium.StrategyName{FcacheL1},
			stg:   pipe([]gopium.Strategy{fcachel1.Curator(b.Curator)}),
		},
		"`full_cache_rounding_cpu_l2` name should return expected strategy": {
			names: []gopium.StrategyName{FcacheL2},
			stg:   pipe([]gopium.Strategy{fcachel2.Curator(b.Curator)}),
		},
		"`full_cache_rounding_cpu_l3` name should return expected strategy": {
			names: []gopium.StrategyName{FcacheL3},
			stg:   pipe([]gopium.Strategy{fcachel3.Curator(b.Curator)}),
		},
		// top, bottom separate pads
		"`separate_padding_system_alignment_top` name should return expected strategy": {
			names: []gopium.StrategyName{SepSysT},
			stg:   pipe([]gopium.Strategy{sepsyst.Curator(b.Curator)}),
		},
		"`separate_padding_cpu_l1_top` name should return expected strategy": {
			names: []gopium.StrategyName{SepL1T},
			stg:   pipe([]gopium.Strategy{sepl1t.Curator(b.Curator)}),
		},
		"`separate_padding_cpu_l2_top` name should return expected strategy": {
			names: []gopium.StrategyName{SepL2T},
			stg:   pipe([]gopium.Strategy{sepl2t.Curator(b.Curator)}),
		},
		"`separate_padding_cpu_l3_top` name should return expected strategy": {
			names: []gopium.StrategyName{SepL3T},
			stg:   pipe([]gopium.Strategy{sepl3t.Curator(b.Curator)}),
		},
		"`separate_padding_system_alignment_bottom` name should return expected strategy": {
			names: []gopium.StrategyName{SepSysB},
			stg:   pipe([]gopium.Strategy{sepsysb.Curator(b.Curator)}),
		},
		"`separate_padding_cpu_l1_bottom` name should return expected strategy": {
			names: []gopium.StrategyName{SepL1B},
			stg:   pipe([]gopium.Strategy{sepl1b.Curator(b.Curator)}),
		},
		"`separate_padding_cpu_l2_bottom` name should return expected strategy": {
			names: []gopium.StrategyName{SepL2B},
			stg:   pipe([]gopium.Strategy{sepl2b.Curator(b.Curator)}),
		},
		"`separate_padding_cpu_l3_bottom` name should return expected strategy": {
			names: []gopium.StrategyName{SepL3B},
			stg:   pipe([]gopium.Strategy{sepl3b.Curator(b.Curator)}),
		},
		// tag processors and modifiers
		"`process_tag_group` name should return expected strategy": {
			names: []gopium.StrategyName{PTGrp},
			stg:   pipe([]gopium.Strategy{ptgrp.Builder(b)}),
		},
		"`add_tag_group_soft` name should return expected strategy": {
			names: []gopium.StrategyName{AddTagS},
			stg:   pipe([]gopium.Strategy{tags.Names(AddTagS)}),
		},
		"`add_tag_group_force` name should return expected strategy": {
			names: []gopium.StrategyName{AddTagF},
			stg:   pipe([]gopium.Strategy{tagf.Names(AddTagF)}),
		},
		"`add_tag_group_discrete` name should return expected strategy": {
			names: []gopium.StrategyName{AddTagSD},
			stg:   pipe([]gopium.Strategy{tagsd.Names(AddTagSD)}),
		},
		"`add_tag_group_force_discrete` name should return expected strategy": {
			names: []gopium.StrategyName{AddTagFD},
			stg:   pipe([]gopium.Strategy{tagfd.Names(AddTagFD)}),
		},
		"`remove_tag_group` name should return expected strategy": {
			names: []gopium.StrategyName{RmTagF},
			stg:   pipe([]gopium.Strategy{tagf}),
		},
		// doc and comment annotations
		"`doc_fields_annotate` name should return expected strategy": {
			names: []gopium.StrategyName{FNoteDoc},
			stg:   pipe([]gopium.Strategy{fnotedoc}),
		},
		"`comment_fields_annotate` name should return expected strategy": {
			names: []gopium.StrategyName{FNoteCom},
			stg:   pipe([]gopium.Strategy{fnotecom}),
		},
		"`doc_struct_annotate` name should return expected strategy": {
			names: []gopium.StrategyName{StNoteDoc},
			stg:   pipe([]gopium.Strategy{stnotedoc}),
		},
		"`comment_struct_annotate` name should return expected strategy": {
			names: []gopium.StrategyName{StNoteCom},
			stg:   pipe([]gopium.Strategy{stnotecom}),
		},
		// lexicographical, length, embedded, exported sorts
		"`name_lexicographical_ascending` name should return expected strategy": {
			names: []gopium.StrategyName{NLexAsc},
			stg:   pipe([]gopium.Strategy{nlexasc}),
		},
		"`name_lexicographical_descending` name should return expected strategy": {
			names: []gopium.StrategyName{NLexDesc},
			stg:   pipe([]gopium.Strategy{nlexdesc}),
		},
		"`type_lexicographical_ascending` name should return expected strategy": {
			names: []gopium.StrategyName{TLexAsc},
			stg:   pipe([]gopium.Strategy{tlexasc}),
		},
		"`type_lexicographical_descending` name should return expected strategy": {
			names: []gopium.StrategyName{TLexDesc},
			stg:   pipe([]gopium.Strategy{tlexdesc}),
		},
		// filters and others
		"`filter_pads` name should return expected strategy": {
			names: []gopium.StrategyName{FPad},
			stg:   pipe([]gopium.Strategy{fpad}),
		},
		"`ignore` name should return expected strategy": {
			names: []gopium.StrategyName{Ignore},
			stg:   pipe([]gopium.Strategy{ignr}),
		},
		"empty name should return empty pipe": {
			stg: pipe{},
		},
		"invalid name should return builder error": {
			names: []gopium.StrategyName{"test"},
			err:   fmt.Errorf(`strategy "test" wasn't found`),
		},
		"complex name should return expected strategy": {
			names: []gopium.StrategyName{Ignore, AddTagS},
			stg:   pipe([]gopium.Strategy{ignr, tags.Names(Ignore, AddTagS)}),
		},
		"invalid name inside complex name should return builder error": {
			names: []gopium.StrategyName{Ignore, "test", AddTagS},
			err:   fmt.Errorf(`strategy "test" wasn't found`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			stg, err := b.Build(tcase.names...)
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
