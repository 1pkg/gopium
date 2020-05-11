package collections

import (
	"fmt"
	"strings"

	"1pkg/gopium"
)

// Hierarchic defines strucs hierarchic collection
// which is categorized by pair of loc and id
type Hierarchic struct {
	rcat string
	cats map[string]Flat
}

// NewHierarchic creates new hierarchic
// collection with root category
func NewHierarchic(rcat string) Hierarchic {
	return Hierarchic{
		rcat: rcat,
		cats: make(map[string]Flat),
	}
}

// Push adds struct to hierarchic collection
func (h Hierarchic) Push(key string, cat string, sts ...gopium.Struct) {
	// remove root cat from the cat
	cat = strings.Replace(cat, h.rcat, "", 1)
	// if cat hasn't been created yet
	flat, ok := h.cats[cat]
	if !ok {
		flat = make(Flat)
	}
	// push not structs to flat collection
	switch l := len(sts); {
	case l == 1:
		flat[key] = sts[0]
	case l > 1:
		// if we have list of struct
		// make unique keys
		for i, st := range sts {
			flat[fmt.Sprintf("%s-%d", key, i)] = st
		}
	}
	// update hierarchic structs collection
	h.cats[cat] = flat
}

// Cat returns hierarchic categoty
// flat collection if any exists
func (h Hierarchic) Cat(cat string) (Flat, bool) {
	// remove root cat from the cat
	cat = strings.Replace(cat, h.rcat, "", 1)
	flat, ok := h.cats[cat]
	return flat, ok
}

// Flat converts hierarchic collection to flat collection
func (h Hierarchic) Flat() Flat {
	// collect all structs by key
	flat := make(Flat)
	for _, lsts := range h.cats {
		for key, st := range lsts {
			flat[key] = st
		}
	}
	return flat
}

// Len calculates total len of hierarchic collection
func (h Hierarchic) Len() int {
	var l int
	for _, cat := range h.cats {
		l += len(cat)
	}
	return l
}
