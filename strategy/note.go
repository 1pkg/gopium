package strategy

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
// that adds size comment annotation
// for all structure fields
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
	var sum, align int64
	for i := range r.Fields {
		f := &r.Fields[i]
		note := gopium.Stamp(fmt.Sprintf("field size: %d align: %d in bytes", f.Size, f.Align))
		if stg.doc {
			f.Comment = append(f.Doc, note)
		} else {
			f.Comment = append(f.Comment, note)
		}
		sum += f.Size
		if align < f.Align {
			align = f.Align
		}
	}
	// note whole structure with size comment
	note := gopium.Stamp(fmt.Sprintf("struct size: %d align: %d in bytes", sum, align))
	if stg.doc {
		r.Doc = append(r.Doc, note)
	} else {
		r.Comment = append(r.Comment, note)
	}
	return
}
