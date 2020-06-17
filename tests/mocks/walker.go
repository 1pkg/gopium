package mocks

import (
	"context"
	"regexp"
	"time"

	"github.com/1pkg/gopium/gopium"
)

// Walker defines mock walker implementation
type Walker struct {
	Err  error         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Wait time.Duration `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_    [8]byte       `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Visit mock implementation
func (w Walker) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	// check error at start
	if w.Err != nil {
		return w.Err
	}
	// sleep for duration if any
	if w.Wait > 0 {
		time.Sleep(w.Wait)
	}
	// return context error otherwise
	return ctx.Err()
}

// WalkerBuilder defines mock walker builder implementation
type WalkerBuilder struct {
	Walker gopium.Walker `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Err    error         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Build mock implementation
func (b WalkerBuilder) Build(gopium.WalkerName) (gopium.Walker, error) {
	return b.Walker, b.Err
}
