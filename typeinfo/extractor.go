package typeinfo

import "go/types"

// Extractor defines typeinfo extractor abstraction
// to fetch typeinfo from provided type
type Extractor interface {
	Extract(types.Type) TypeInfo
}

// ExtractorTypesSizes defines typeinfo Extractor default types.Sizes implementation
// that uses Sizes.Sizeof in order to get size of the type
type ExtractorTypesSizes struct {
	sizes types.Sizes
}

// NewExtractorTypesSizes creates instance of typeinfo ExtractorTypesSizes
// and requires compiler and arch for types.Sizes initialization
func NewExtractorTypesSizes(compiler, arch string) ExtractorTypesSizes {
	return ExtractorTypesSizes{sizes: types.SizesFor(compiler, arch)}
}

// Extract typeinfo Extractor default types.Sizes implementation
func (ets ExtractorTypesSizes) Extract(t types.Type) TypeInfo {
	size := ets.sizes.Sizeof(t)
	return TypeInfo{Type: t.String(), Size: size}
}
