package pkgs_types

import "go/types"

// WhistleblowerGoTypes defines whistleblower default "go/types" implementation
// that uses types.Sizes Sizeof in order to get name and size from provided type
type WhistleblowerGoTypes struct {
	sizes types.Sizes
}

// NewWhistleblowerGoTypes creates instance of ExtractorGoTypes
// and requires compiler and arch for types.Sizes initialization
func NewWhistleblowerGoTypes(compiler, arch string) WhistleblowerGoTypes {
	return WhistleblowerGoTypes{sizes: types.SizesFor(compiler, arch)}
}

// Expose WhistleblowerGoTypes implementation
func (wb WhistleblowerGoTypes) Expose(t types.Type) (name string, size int64) {
	name = t.String()
	size = wb.sizes.Sizeof(t)
	return
}
