package mocks

import "context"

// Runner defines mock runner implementation
type Runner struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Run mock implementation
func (r Runner) Run(context.Context) error {
	return r.Err
}
