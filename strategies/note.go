package strategies

import (
	"context"
	"fmt"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of note presets
var (
	fnotedoc  = note{doc: true, field: true}
	fnotecom  = note{doc: false, field: true}
	stnotedoc = note{doc: true, field: false}
	stnotecom = note{doc: false, field: false}
)

// note defines strategy implementation
// that adds size doc or comment annotation
// for each structure field
// and aggregated size annotation for structure
type note struct {
	doc   bool `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	field bool `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 2 bytes; struct align: 1 bytes; struct aligned size: 2 bytes; - 🌺 gopium @1pkg

// Apply note implementation
func (stg note) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// preset defaults
	var size, alsize, align int64 = 0, 0, 1
	// prepare fields slice
	if flen := len(r.Fields); flen > 0 {
		// note each field with size comment
		rfields := make([]gopium.Field, 0, flen)
		collections.WalkStruct(r, 0, func(pad int64, fields ...gopium.Field) {
			// add pad to aligned size
			alsize += pad
			for _, f := range fields {
				// note only in field mode
				if stg.field {
					// create note comment
					note := fmt.Sprintf(
						"// field size: %d bytes; field align: %d bytes; - %s",
						f.Size,
						f.Align,
						gopium.STAMP,
					)
					if stg.doc {
						f.Doc = append(f.Doc, note)
					} else {
						f.Comment = append(f.Comment, note)
					}
				}
				// add field size to both sizes
				size += f.Size
				alsize += f.Size
				// update struct align size
				// if field align size is bigger
				if f.Align > align {
					align = f.Align
				}
				// append field to result
				rfields = append(rfields, f)
			}
		})
		// update result fields
		r.Fields = rfields
	}
	// note structure with size comment
	// note only in non field mode
	if !stg.field {
		// create note comment
		note := fmt.Sprintf(
			"// struct size: %d bytes; struct align: %d bytes; struct aligned size: %d bytes; - %s",
			size,
			align,
			alsize,
			gopium.STAMP,
		)
		if stg.doc {
			r.Doc = append(r.Doc, note)
		} else {
			r.Comment = append(r.Comment, note)
		}
	}
	return r, ctx.Err()
}
