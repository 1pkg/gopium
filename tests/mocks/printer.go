package mocks

import (
	"context"
	"go/ast"
	"go/token"
	"io"
)

// Printer defines mock ast printer implementation
type Printer struct {
	Err error `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Print mock implementation
func (p Printer) Print(context.Context, io.Writer, *token.FileSet, ast.Node) error {
	return p.Err
}
