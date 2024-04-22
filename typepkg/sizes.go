package typepkg

import (
	"go/types"

	"github.com/1pkg/gopium/collections"
)

// ftype is a stub to extract WordSize and MaxAlign out of types.Sizes interface across gc and gccgo compilers.
type stubtype struct{}

func (t stubtype) Underlying() types.Type {
	return stubtype{}
}

func (t stubtype) String() string {
	return ""
}

// stdsizes implements sizes interface using types std sizes
type stdsizes struct {
	types.Sizes `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 8 bytes; struct align: 8 bytes; struct aligned size: 8 bytes; struct ptr scan size: 8 bytes; - 🌺 gopium @1pkg

func (s stdsizes) WordSize() int64 {
	// This should work for both gc and gccgo types as default case returning proper WordSize.
	return s.Sizeof(stubtype{})
}

func (s stdsizes) MaxAlign() int64 {
	// This should work for both gc and gccgo types as default case returning proper MaxAlign.
	return s.Alignof(types.Typ[types.Complex128])
}

// Ptr implementation is vendored from
// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.9.3:go/analysis/passes/fieldalignment/fieldalignment.go;l=330
func (s stdsizes) Ptr(t types.Type) int64 {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String, types.UnsafePointer:
			return s.WordSize()
		}
		return 0
	case *types.Chan, *types.Map, *types.Pointer, *types.Signature, *types.Slice:
		return s.WordSize()
	case *types.Interface:
		return 2 * s.WordSize()
	case *types.Array:
		n := t.Len()
		if n == 0 {
			return 0
		}
		a := s.Ptr(t.Elem())
		if a == 0 {
			return 0
		}
		z := s.Sizeof(t.Elem())
		return (n-1)*z + a
	case *types.Struct:
		nf := t.NumFields()
		if nf == 0 {
			return 0
		}

		var o, p int64
		for i := 0; i < nf; i++ {
			ft := t.Field(i).Type()
			a, sz := s.Alignof(ft), s.Sizeof(ft)
			fp := s.Ptr(ft)
			o = collections.Align(o, a)
			if fp != 0 {
				p = o + fp
			}
			o += sz
		}
		return p
	}

	return 0
}
