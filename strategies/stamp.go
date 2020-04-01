package strategies

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of stamp presets
var (
	stampdoc = stamp{doc: true}
	stampcom = stamp{doc: false}
)

// stamp defines strategy implementation
// that adds doc `auto curated` stamp to structure doc
type stamp struct {
	doc bool
}

// Apply stamp implementation
func (stg stamp) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// create stamp
	stamp := fmt.Sprintf(
		"// struct %q has been auto curated by gopium - %s",
		r.Name,
		gopium.STAMP,
	)
	// add stamp to structure doc or comment
	if stg.doc {
		r.Doc = append(r.Doc, stamp)
	} else {
		r.Comment = append(r.Comment, stamp)
	}
	return
}
