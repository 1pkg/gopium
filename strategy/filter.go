package strategy

import (
	"context"
	"regexp"

	"1pkg/gopium"
)

// list of filter presets
var (
	filterpad = filter{regexp.MustCompile(`^_$`)}
)

// filter defines strategy implementation
// that filters all fields
// that match provided regex
type filter struct {
	regex *regexp.Regexp
}

// Apply filter implementation
func (stg filter) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// prepare filtered fields list
	fields := make([]gopium.Field, 0, len(r.Fields))
	// then go though all original fields
	for _, f := range r.Fields {
		// check if field name matches regex
		if stg.regex.MatchString(r.Name) {
			continue
		}
		// if it doesn't append it to fields
		fields = append(fields, f)
	}
	// update result field list
	r.Fields = fields
	return
}
