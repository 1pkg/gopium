package strategies

import (
	"context"
	"fmt"

	"1pkg/gopium"
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
// and aggregated size annotation
// for whole structure
type note struct {
	doc   bool
	field bool
}

// Apply note implementation
func (stg note) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o.Copy()
	// note each field with size comment
	var size, align int64
	for i := range r.Fields {
		f := &r.Fields[i]
		// note only in field mode
		if stg.field {
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
		// calculate total size and align
		size += f.Size
		if align < f.Align {
			align = f.Align
		}
	}
	// note whole structure with size comment
	// note only in non field mode
	if !stg.field {
		note := fmt.Sprintf(
			"// struct size: %d bytes; struct align: %d bytes; - %s",
			size,
			align,
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
