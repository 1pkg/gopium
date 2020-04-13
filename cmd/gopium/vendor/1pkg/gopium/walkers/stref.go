package walkers

import (
	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// stRef helps to create gopium.Struct
// size refence for provided key
// by preallocating the key and then
// pushing total struct size to ref with closure
func stRef(ref *collections.Reference, name string) func(gopium.Struct) {
	// preallocate the key
	ref.Alloc(name)
	// return the pushing closure
	return func(st gopium.Struct) {
		// calculate total struct size
		var size int64
		for _, f := range st.Fields {
			size += f.Size
		}
		// set ref key size
		ref.Set(name, size)
	}
}
