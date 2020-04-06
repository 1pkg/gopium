package strategies

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of note presets
var (
	notedoc = note{doc: true}
	notecom = note{doc: false}
)

// note defines strategy implementation
// that adds size doc or comment annotation
// for each structure field
// and aggregated size annotation
// for whole structure
type note struct {
	doc bool
}

// Apply note implementation
func (stg note) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// note each field with size comment
	var size, align int64
	for i := range r.Fields {
		f := &r.Fields[i]
		note := fmt.Sprintf(
			"// field size: %d bytes field align: %d bytes - %s",
			f.Size,
			f.Align,
			gopium.STAMP,
		)
		if stg.doc {
			f.Doc = append(f.Doc, note)
		} else {
			f.Comment = append(f.Comment, note)
		}
		size += f.Size
		if align < f.Align {
			align = f.Align
		}
	}
	// note whole structure with size comment
	note := fmt.Sprintf(
		"// struct size: %d bytes struct align: %d bytes - %s",
		size,
		align,
		gopium.STAMP,
	)
	if stg.doc {
		r.Doc = append(r.Doc, note)
	} else {
		r.Comment = append(r.Comment, note)
	}
	return
}
