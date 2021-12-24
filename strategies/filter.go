package strategies

import (
	"context"
	"regexp"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of filter presets
var (
	// list of filter presets
	fpad = filter{
		nregex: regexp.MustCompile(`^_$`),
	}
)

// filter defines strategy implementation
// that filters out all structure fields
// that matches provided criteria
type filter struct {
	nregex *regexp.Regexp `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	tregex *regexp.Regexp `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; struct ptr scan size: 16 bytes; - 🌺 gopium @1pkg

// Apply filter implementation
func (stg filter) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// prepare filtered fields slice
	if flen := len(r.Fields); flen > 0 {
		fields := make([]gopium.Field, 0, flen)
		// then go though all original fields
		for _, f := range r.Fields {
			// check if field name matches regex
			if stg.nregex != nil && stg.nregex.MatchString(f.Name) {
				continue
			}
			// check if field type matches regex
			if stg.tregex != nil && stg.tregex.MatchString(f.Type) {
				continue
			}
			// if it doesn't append it to fields
			fields = append(fields, f)
		}
		// update result fields
		r.Fields = fields
	}
	return r, ctx.Err()
}
