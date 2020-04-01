package strategy

import (
	"context"
	"regexp"

	"1pkg/gopium"
)

// list of filter presets
var (
	// to make bools them addressable
	t = true
	f = false
	// list of filter presets
	fpad = filter{
		nregex: regexp.MustCompile(`^_$`),
	}
	femb = filter{
		emb: &t,
	}
	fnotemb = filter{
		emb: &f,
	}
	fexp = filter{
		exp: &t,
	}
	fnotexp = filter{
		exp: &f,
	}
)

// filter defines strategy implementation
// that filters all fields
// that match provided regex
type filter struct {
	nregex, tregex *regexp.Regexp
	emb, exp       *bool
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
		if stg.nregex != nil && stg.nregex.MatchString(f.Name) {
			continue
		}
		// check if field type matches regex
		if stg.tregex != nil && stg.tregex.MatchString(f.Type) {
			continue
		}
		// check if field embedded matches condition
		if stg.emb != nil && *stg.emb == f.Embedded {
			continue
		}
		// check if field exported matches condition
		if stg.exp != nil && *stg.exp == f.Exported {
			continue
		}
		// if it doesn't append it to fields
		fields = append(fields, f)
	}
	// update result field list
	r.Fields = fields
	return
}
