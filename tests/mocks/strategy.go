package mocks

import (
	"context"

	"github.com/1pkg/gopium/gopium"
)

// Strategy defines mock strategy implementation
type Strategy struct {
	R   gopium.Struct `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Err error         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_   [24]byte      `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 128 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; - 🌺 gopium @1pkg

// Apply mock implementation
func (stg *Strategy) Apply(context.Context, gopium.Struct) (gopium.Struct, error) {
	return stg.R, stg.Err
}

// StrategyBuilder defines mock strategy builder implementation
type StrategyBuilder struct {
	Strategy gopium.Strategy `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Err      error           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Build mock implementation
func (b StrategyBuilder) Build(...gopium.StrategyName) (gopium.Strategy, error) {
	return b.Strategy, b.Err
}
