package strategies

import (
	"context"
	"sort"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of pack presets
var (
	pck = pack{}
)

// pack defines strategy implementation
// that rearranges structure fields
// to obtain optimal memory utilization
// by sorting fields accordingly
// to their aligns and sizes in some order
type pack struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Apply pack implementation
func (stg pack) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// execute memory sorting
	// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.1.7:go/analysis/passes/fieldalignment/fieldalignment.go;l=145;bpv=0;bpt=1
	sort.SliceStable(r.Fields, func(i, j int) bool {
		// place zero sized objects before non-zero sized objects
		zeroi, zeroj := r.Fields[i].Size == 0, r.Fields[j].Size == 0
		if zeroi != zeroj {
			return zeroi
		}
		// then compare aligns of two fields
		// bigger aligmnet means upper position
		if r.Fields[i].Align != r.Fields[j].Align {
			return r.Fields[i].Align > r.Fields[j].Align
		}
		// place pointerful objects before pointer-free objects
		noptri, noptrj := r.Fields[i].Ptr == 0, r.Fields[j].Ptr == 0
		if noptri != noptrj {
			return noptrj
		}

		if !noptri {
			// if both have pointers
			// then place objects with less trailing
			// non-pointer bytes earlier;
			// that is, place the field with the most trailing
			// non-pointer bytes at the end of the pointerful section
			traili, trailj := r.Fields[i].Size-r.Fields[i].Ptr, r.Fields[j].Size-r.Fields[j].Ptr
			if traili != trailj {
				return traili > trailj
			}
		}

		// then compare sizes of two fields
		// bigger size means upper position
		return r.Fields[i].Size > r.Fields[j].Size
	})
	return r, ctx.Err()
}
