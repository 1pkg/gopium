package data

import (
	"io"

	"1pkg/gopium/gopium"
)

// Writer defines tests data writter implementation
// which reuses underlying locator
// but purifies location generation
type Writer struct {
	Writer gopium.CategoryWriter `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// Generate writer implementation
func (w Writer) Generate(loc string) (io.WriteCloser, error) {
	// purify the loc then
	// just reuse underlying writer
	return w.Writer.Generate(purify(loc))
}

// Category writer implementation
func (w Writer) Category(cat string) error {
	return w.Writer.Category(cat)
}
