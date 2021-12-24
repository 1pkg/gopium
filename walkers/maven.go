package walkers

import (
	"go/types"
	"sync"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// ptrsizealign defines data transfer
// object that holds type triplet
// of ptr, size and align vals
type ptrsizealign struct {
	ptr   int64   `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	size  int64   `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	align int64   `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_     [8]byte `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg

// maven defines visiting helper
// that aggregates some useful
// operations on underlying facilities
type maven struct {
	exp   gopium.Exposer         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	loc   gopium.Locator         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	ref   *collections.Reference `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	store sync.Map               `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_     [48]byte               `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 128 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; struct ptr scan size: 72 bytes; - 🌺 gopium @1pkg

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
		sa := m.refpsa(f.Type())
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     m.exp.Name(f.Type()),
			Size:     sa.size,
			Align:    sa.align,
			Ptr:      sa.ptr,
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return r
}

// refsa defines ptr and size and align getter
// with reference helper that uses reference
// if it has been provided
// or uses exposer to expose type size
func (m *maven) refpsa(t types.Type) ptrsizealign {
	// in case we don't have a reference
	// just use default exposer size
	if m.ref == nil {
		return ptrsizealign{
			ptr:   m.exp.Ptr(t),
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
			return ptrsizealign{}
		}
		// n > 0
		sa := m.refpsa(tp.Elem())
		sa.size = collections.Align(sa.size, sa.align)*(n-1) + sa.size
		sa.ptr = sa.ptr * n
		return sa
	case *types.Named:
		// in case it's not a struct skip it
		if _, ok := tp.Underlying().(*types.Struct); !ok {
			break
		}
		// get id for named structures
		id := m.loc.ID(tp.Obj().Pos())
		// get size of the structure from ref
		if sa, ok := m.ref.Get(id).(ptrsizealign); ok {
			return sa
		}
	}
	// just use default exposer size
	return ptrsizealign{
		ptr:   m.exp.Ptr(t),
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
		// calculate structure align, aligned size and ptr size
		stsize, stalign, ptrsize := collections.SizeAlignPtr(st)
		// set ref key size and align
		m.ref.Set(name, ptrsizealign{size: stsize, align: stalign, ptr: ptrsize})
	}
}
