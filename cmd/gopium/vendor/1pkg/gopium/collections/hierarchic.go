package collections

import "1pkg/gopium"

// Hierarchic defines strucs hierarchical collection
// which is categorized by pair of loc and id
type Hierarchic map[string]Flat

// Push adds struct to hierarchic collection
func (h Hierarchic) Push(key, cat string, st gopium.Struct) {
	// if loc hasn't been created yet
	flat, ok := h[cat]
	if !ok {
		flat = make(Flat)
	}
	// push struct to flat collection
	flat[key] = st
	// update hierarchic structs collection
	h[cat] = flat
}

// Cat returns hierarchic categoty
// flat collection if any exists
func (h Hierarchic) Cat(loc string) (Flat, bool) {
	flat, ok := h[loc]
	return flat, ok
}

// Flat converts hierarchic collection to flat collection
func (h Hierarchic) Flat() Flat {
	// collect all structs by key
	flat := make(Flat)
	for _, lsts := range h {
		for key, st := range lsts {
			flat[key] = st
		}
	}
	return flat
}
