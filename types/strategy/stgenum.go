package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
	gtypes "1pkg/gopium/types"
)

// enumerate defines struct enumerating strategy implementation
// that goes through all structure fields and uses gtypes.Extractor
// to extract gopium.Field DTO for each field
// and put it back to resulted gopium.Struct object
type stgenum struct {
	extractor gtypes.Extractor
}

// Apply enumerate implementation
func (stg stgenum) Apply(ctx context.Context, name string, st *types.Struct) (r gopium.StructError) {
	// build full hierarchical name of the structure
	r.Struct.Name = fmt.Sprintf("%s/%s", name, st)
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Struct.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get tag
		tag := st.Tag(i)
		// get typeinfo
		tname, tsize := stg.extractor.Extract(f.Type())
		// fill field structure
		r.Struct.Fields = append(r.Struct.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     tname,
			Size:     tsize,
			Tag:      tag,
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return
}
