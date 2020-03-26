package walker

import (
	"go/types"
	"strings"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/walker/reference"
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
func (m *maven) enum(name string, st *types.Struct, ref *reference.Ref) (r gopium.Struct) {
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
func (m *maven) refsize(t types.Type, ref *reference.Ref) int64 {
	// in case we have reference
	if ref != nil {
		// for refsize only structures
		// and arrays should be calculated
		// not with default exposer size
		switch tp := t.(type) {
		case *types.Array:
			// ignore not struct arrays
			if _, ok := tp.Elem().(*types.Struct); !ok {
				break
			}
			// note: copied from `go/types/sizes.go`
			n := tp.Len()
			if n <= 0 {
				return 0
			}
			// n > 0
			a := m.exposer.Align(tp.Elem())
			z := m.refsize(tp.Elem(), ref)
			return gopium.Align(z, a)*(n-1) + z
		case *types.Struct:
			name := tp.String()
			// ignore structs from different pkg
			if strings.Contains(name, ".") {
				break
			}
			// ignore anonymus structs
			if strings.Contains(name, "struct") {
				break
			}
			// get size of the structure from ref
			if size := ref.Get(name); size >= 0 {
				return size
			}
		}
	}
	// just use default exposer size
	return m.exposer.Size(t)
}
