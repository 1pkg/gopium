package types

import "go/types"

// Extractor defines type info extractor abstraction
// to extract name and size from provided type
type Extractor interface {
	Extract(t types.Type) (name string, size int64)
}

// ExtractorGoTypes defines Extractor default "go/types" implementation
// that uses types.Sizes Sizeof in order to get name and size from provided type
type ExtractorGoTypes struct {
	sizes types.Sizes
}

// NewExtractorGoTypes creates instance of ExtractorGoTypes
// and requires compiler and arch for types.Sizes initialization
func NewExtractorGoTypes(compiler, arch string) ExtractorGoTypes {
	return ExtractorGoTypes{sizes: types.SizesFor(compiler, arch)}
}

// Extract ExtractorGoTypes implementation
func (ets ExtractorGoTypes) Extract(t types.Type) (name string, size int64) {
	name = t.String()
	size = ets.sizes.Sizeof(t)
	return
}
