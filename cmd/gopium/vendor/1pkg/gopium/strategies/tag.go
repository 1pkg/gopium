package strategies

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"1pkg/gopium"
)

// gopium tag name
const tagn = "gopium"

// list of tag presets
var (
	tagrm = tag{force: true}
)

// tag defines strategy implementation
// that adds or updates fields tags annotation
// that could be processed by group strategy
type tag struct {
	tag   string
	force bool
}

// Apply tag implementation
func (stg tag) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// iterate through all fields
	for i := range r.Fields {
		f := &r.Fields[i]
		// grab the field tag
		tag, ok := reflect.StructTag(f.Tag).Lookup(tagn)
		// in case gopium tag already exists
		// and force is set - replace tag
		// in case tag is not empty and
		// gopium tag doesn't exist - append tag
		// in case tag is empty - set tag
		ntag := fmt.Sprintf(`%s:"%s"`, tagn, stg.tag)
		if ok && stg.force {
			f.Tag = strings.Replace(f.Tag, tag, stg.tag, 1)
		} else if len(f.Tag) != 0 {
			f.Tag += " " + ntag
		} else {
			f.Tag = ntag
		}
	}
	return
}
