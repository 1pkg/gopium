package walker

import (
	"go/types"
	"sync"

	"1pkg/gopium"
)

// maven defines visiting helper
// that aggregates has and enum facilities
type maven struct {
	exposer gopium.Exposer
	idfunc  gopium.IDFunc
	store   sync.Map
}

// has defines struct store id helper
// that uses gopium.IDFunc to build id
// for a structure and check that
// builded id has not been stored already
func (m *maven) has(tn *types.TypeName) (string, bool) {
	// build id for the structure
	id := m.idfunc(tn.Pos())
	// in case id of structure
	// has been already stored
	if _, ok := m.store.Load(id); ok {
		return id, true
	}
	// mark id of structure as stored
	m.store.Store(id, struct{}{})
	return id, false
}

// enum defines struct enumerating converting helper
// that goes through all structure fields
// and uses gopium.Exposer to expose gopium.Field DTO
// for each field and puts them back
// to resulted gopium.Struct object
func (m *maven) enum(name string, st *types.Struct) (r gopium.Struct) {
	// set structure name
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     m.exposer.Name(f.Type()),
			Size:     m.exposer.Size(f.Type()),
			Align:    m.exposer.Align(f.Type()),
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return
}
