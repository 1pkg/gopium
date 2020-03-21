package strategy

import (
	"context"
	"go/types"

	"1pkg/gopium"
)

// enum defines struct enumerating strategy implementation
// that goes through all structure fields and uses gopium.Maven
// to expose gopium.Field DTO for each field
// and puts it back to resulted gopium.Struct object
type enum struct {
	m gopium.Maven
}

// Apply enum implementation
func (stg enum) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// set structure name
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get tag
		tag := st.Tag(i)
		// get typeinfo
		tname := stg.m.Name(f.Type())
		tsize := stg.m.Size(f.Type())
		talign := stg.m.Align(f.Type())
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     tname,
			Size:     tsize,
			Align:    talign,
			Tag:      tag,
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	o = r
	return
}
