package walkers

import (
	"go/types"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/walkers/ref"
)

// maven defines visiting helper
// that aggregates has and enum facilities
type maven struct {
	exposer gopium.Exposer
	locator gopium.Locator
	store   sync.Map
}

// has defines struct store id helper
// that uses gopium.Locator to build id
// for a structure and check that
// builded id has not been stored already
func (m *maven) has(tn *types.TypeName) (id, loc string, ok bool) {
	// build id for the structure
	id = m.locator.ID(tn.Pos())
	// build loc for the structure
	loc = m.locator.Loc(tn.Pos())
	// in case id of structure
	// has been already stored
	if _, ok := m.store.Load(id); ok {
		return id, loc, true
	}
	// mark id of structure as stored
	m.store.Store(id, struct{}{})
	return id, loc, false
}

// enum defines struct enumerating converting helper
// that goes through all structure fields
// and uses gopium.Exposer to expose gopium.Field DTO
// for each field and puts them back
// to resulted gopium.Struct object
func (m *maven) enum(name string, st *types.Struct, ref *ref.Ref) (r gopium.Struct) {
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
			Size:     m.refsize(f.Type(), ref),
			Align:    m.exposer.Align(f.Type()),
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return
}

// refsize defines size getter with reference helper
// that uses reference if it has been provided
// or uses gopium.Exposer to expose type size
func (m *maven) refsize(t types.Type, ref *ref.Ref) int64 {
	// in case we have reference
	if ref != nil {
		// for refsize only named structures
		// and arrays should be calculated
		// not with default exposer size
		switch tp := t.(type) {
		case *types.Array:
			// note: copied from `go/types/sizes.go`
			n := tp.Len()
			if n <= 0 {
				return 0
			}
			// n > 0
			a := m.exposer.Align(tp.Elem())
			z := m.refsize(tp.Elem(), ref)
			return gopium.Align(z, a)*(n-1) + z
		case *types.Named:
			// in case it's not a struct skip it
			if _, ok := tp.Underlying().(*types.Struct); ok {
				break
			}
			// get id for named structures
			id := m.locator.ID(tp.Obj().Pos())
			// get size of the structure from ref
			if size := ref.Get(id); size >= 0 {
				return size
			}
		}
	}
	// just use default exposer size
	return m.exposer.Size(t)
}
