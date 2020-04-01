package ref

import "1pkg/gopium"

// StRef helps to create gopium.Struct
// size refence for provided key
// by preallocating the key and then
// pushing total struct size to ref with closure
func (r *Ref) StRef(key string) func(gopium.Struct) {
	// preallocate the key
	r.alloc(key)
	// return the pushing closure
	return func(st gopium.Struct) {
		// calculate total struct size
		var size int64
		for _, f := range st.Fields {
			size += f.Size
		}
		// set ref key size
		r.Set(key, size)
	}
}