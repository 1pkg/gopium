package strategies

import (
	"context"
	"fmt"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// list of stamp presets
var (
	stampdoc = stamp{doc: true}
	stampcom = stamp{doc: false}
)

// stamp defines strategy implementation
// that adds doc or comment gopium stamp to structure
type stamp struct {
	doc bool
}

// Apply stamp implementation
func (stg stamp) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// create stamp
	stamp := fmt.Sprintf(
		"// struct has been auto curated - %s",
		gopium.STAMP,
	)
	// add stamp to structure doc or comment
	if stg.doc {
		r.Doc = append(r.Doc, stamp)
	} else {
		r.Comment = append(r.Comment, stamp)
	}
	return r, ctx.Err()
}
