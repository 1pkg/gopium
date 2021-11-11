package typepkg

import (
	"go/types"

	"github.com/1pkg/gopium/collections"
)

// sizes is types sizes plus ptr data size interface
type sizes interface {
	types.Sizes
	Ptr(t types.Type) int64
}

// stdsizes implements sizes interace using types std sizes
type stdsizes struct {
	*types.StdSizes
}

// Ptr implementation is vendored from
// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.1.7:go/analysis/passes/fieldalignment/fieldalignment.go;l=324
func (s stdsizes) Ptr(t types.Type) int64 {
	switch t := t.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String, types.UnsafePointer:
			return s.WordSize
		}
		return 0
	case *types.Chan, *types.Map, *types.Pointer, *types.Signature, *types.Slice:
		return s.WordSize
	case *types.Interface:
		return 2 * s.WordSize
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
	panic("impossible")
}
