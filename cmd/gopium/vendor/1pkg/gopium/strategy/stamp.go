package strategy

import (
	"context"
	"fmt"

	"1pkg/gopium"
)

// list of stamp presets
var (
	stmp = stamp{}
)

// stamp defines strategy implementation
// that adds doc `auto curated` stamp to structure doc
type stamp struct{}

// Apply stamp implementation
func (stg stamp) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// create stamp
	stamp := fmt.Sprintf("struct %q has been auto curated", r.Name)
	stamp = gopium.Stamp(stamp)
	// add stamp to structure doc
	r.Doc = append(r.Doc, stamp)
	return
}
