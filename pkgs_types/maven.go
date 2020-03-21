package pkgs_types

import "go/types"

// MavenGoTypes defines maven default "go/types" implementation
// that uses types.Sizes Sizeof in order to get type info
type MavenGoTypes struct {
	sizes types.Sizes
}

// NewWhistleblowerGoTypes creates instance of ExtractorGoTypes
// and requires compiler and arch for types.Sizes initialization
func NewMavenGoTypes(compiler, arch string) MavenGoTypes {
	return MavenGoTypes{sizes: types.SizesFor(compiler, arch)}
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
	// TODO
	return -1
}

// Name MavenGoTypes implementation
func (m MavenGoTypes) Name(t types.Type) string {
	return t.String()
}

// Size MavenGoTypes implementation
func (m MavenGoTypes) Size(t types.Type) int64 {
	return m.sizes.Sizeof(t)
}

// Size MavenGoTypes implementation
func (m MavenGoTypes) Align(t types.Type) int64 {
	return m.sizes.Alignof(t)
}
