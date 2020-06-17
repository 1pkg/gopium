package mocks

import (
	"context"
	"go/ast"

	"github.com/1pkg/gopium/gopium"
)

// Walk defines mock ast walker implementation
type Walk struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Walk mock implementation
func (w Walk) Walk(context.Context, ast.Node, gopium.Visitor, gopium.Comparator) (ast.Node, error) {
	return nil, w.Err
}
