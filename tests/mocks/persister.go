package mocks

import (
	"context"
	"go/ast"

	"1pkg/gopium/gopium"
)

// Persister defines mock pesister implementation
type Persister struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Persist mock implementation
func (p Persister) Persist(context.Context, gopium.Printer, gopium.Writer, gopium.Locator, ast.Node) error {
	return p.Err
}
