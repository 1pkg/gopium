package strategy

import (
	"context"
	"go/types"
	"regexp"

	"1pkg/gopium"
)

// filter defines struct fields filter strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then filters all fields that match provided regex
type filter struct {
	m gopium.Maven
	r *regexp.Regexp
}

// Apply filter implementation
func (stg filter) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// prepare filtred fields list
	fields := make([]gopium.Field, 0, len(r.Fields))
	// then execute memory sorting
	for _, f := range r.Fields {
		// check if field name matches regex
		if stg.r.MatchString(name) {
			continue
		}
		// if it doesn't append it to fields
		fields = append(fields, f)
	}
	// update result field list
	r.Fields = fields
	return
}
