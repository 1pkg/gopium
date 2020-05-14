package strategies

import (
	"context"
	"regexp"

	"1pkg/gopium"
)

// list of filter presets
var (
	// to make bools addressable
	tvar = true
	fvar = false
	// list of filter presets
	fpad = filter{
		nregex: regexp.MustCompile(`^_$`),
	}
	femb = filter{
		emb: &tvar,
	}
	fnotemb = filter{
		emb: &fvar,
	}
	fexp = filter{
		exp: &tvar,
	}
	fnotexp = filter{
		exp: &fvar,
	}
)

// filter defines strategy implementation
// that filters out all structure fields
// that matches provided criteria
type filter struct {
	nregex, tregex *regexp.Regexp
	emb, exp       *bool
}

// Apply filter implementation
func (stg filter) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o.Copy()
	// prepare filtered fields slice
	if flen := len(r.Fields); flen > 0 {
		fields := make([]gopium.Field, 0, flen)
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
		// update result fields
		r.Fields = fields
	}
	return r, ctx.Err()
}
