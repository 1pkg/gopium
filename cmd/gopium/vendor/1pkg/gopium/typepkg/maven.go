package typepkg

import "go/types"

// MavenGoTypes defines maven default "go/types" implementation
// that uses types.Sizes Sizeof in order to get type info
type MavenGoTypes struct {
	sizes  types.Sizes
	caches map[uint]int64
}

// NewWhistleblowerGoTypes creates instance of ExtractorGoTypes
// and requires compiler and arch for types.Sizes initialization
func NewMavenGoTypes(compiler, arch string, caches ...int64) MavenGoTypes {
	// go through all passed caches
	// and fill them to cache map
	cm := make(map[uint]int64, len(caches))
	for i, cache := range caches {
		cm[uint(i+1)] = cache
	}
	return MavenGoTypes{
		sizes:  types.SizesFor(compiler, arch),
		caches: cm,
	}
}

// SysWord MavenGoTypes implementation
func (m MavenGoTypes) SysWord() int64 {
	return m.sizes.(*types.StdSizes).WordSize
}

// SysAlign MavenGoTypes implementation
func (m MavenGoTypes) SysAlign() int64 {
	return m.sizes.(*types.StdSizes).MaxAlign
}

// SysCache MavenGoTypes implementation
func (m MavenGoTypes) SysCache(level uint) int64 {
	// if we have specified cache size
	if size, ok := m.caches[level]; ok {
		return size
	}
	// otherwise just return
	// typical cpu cache size
	return 64
}

// Name MavenGoTypes implementation
func (m MavenGoTypes) Name(t types.Type) string {
	return t.String()
}

// Size MavenGoTypes implementation
func (m MavenGoTypes) Size(t types.Type) int64 {
	return m.sizes.Sizeof(t)
}

// Align MavenGoTypes implementation
func (m MavenGoTypes) Align(t types.Type) int64 {
	return m.sizes.Alignof(t)
}
