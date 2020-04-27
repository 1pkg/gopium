package walkers

import (
	"go/types"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// sizealign defines data transfer
// object that holds type pair
// of size and align vals
type sizealign struct {
	size  int64
	align int64
}

// maven defines visiting helper
// that aggregates some useful
// operations on underlying facilities
type maven struct {
	exp   gopium.Exposer
	loc   gopium.Locator
	store sync.Map
	ref   *collections.Reference
}

// has defines struct store id helper
// that uses locator to build id
// for a structure and check that
// builded id has not been stored already
func (m *maven) has(tn *types.TypeName) (id string, loc string, ok bool) {
	// build id for the structure
	id = m.loc.ID(tn.Pos())
	// build loc for the structure
	loc = m.loc.Loc(tn.Pos())
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
// and uses exposer to expose field DTO
// for each field and puts them back
// to resulted struct object
func (m *maven) enum(name string, st *types.Struct) gopium.Struct {
	// set structure name
	r := gopium.Struct{}
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get size and align for field
		sa := m.refsa(f.Type())
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     m.exp.Name(f.Type()),
			Size:     sa.size,
			Align:    sa.align,
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return r
}

// refsa defines size and align getter
// with reference helper that uses reference
// if it has been provided
// or uses exposer to expose type size
func (m *maven) refsa(t types.Type) sizealign {
	// in case we don't have a reference
	// just use default exposer size
	if m.ref == nil {
		return sizealign{
			size:  m.exp.Size(t),
			align: m.exp.Align(t),
		}
	}
	// for refsize only named structures
	// and arrays should be calculated
	// not with default exposer size
	switch tp := t.(type) {
	case *types.Array:
		// note: copied from `go/types/sizes.go`
		n := tp.Len()
		if n <= 0 {
			return sizealign{}
		}
		// n > 0
		sa := m.refsa(tp.Elem())
		sa.size = gopium.Align(sa.size, sa.align)*(n-1) + sa.size
		return sa
	case *types.Named:
		// in case it's not a struct skip it
		if _, ok := tp.Underlying().(*types.Struct); !ok {
			break
		}
		// get id for named structures
		id := m.loc.ID(tp.Obj().Pos())
		// get size of the structure from ref
		if sa, ok := m.ref.Get(id).(sizealign); ok {
			return sa
		}
	}
	// just use default exposer size
	return sizealign{
		size:  m.exp.Size(t),
		align: m.exp.Align(t),
	}
}

// refst helps to create struct
// size refence for provided key
// by preallocating the key and then
// pushing total struct size to ref with closure
func (m *maven) refst(name string) func(gopium.Struct) {
	// preallocate the key
	m.ref.Alloc(name)
	// return the pushing closure
	return func(st gopium.Struct) {
		// calculate total struct size
		var size, align int64 = 0, 1
		for _, f := range st.Fields {
			if f.Align > align {
				align = f.Align
			}
			size += f.Size
		}
		// set ref key size and align
		m.ref.Set(name, sizealign{size: size, align: align})
	}
}
