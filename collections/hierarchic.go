package collections

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/1pkg/gopium/gopium"
)

// Hierarchic defines strucs hierarchic collection
// which is categorized by pair of loc and id
type Hierarchic struct {
	rcat string          `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	cats map[string]Flat `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
	_    [8]byte         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,comment_struct_annotate,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

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
func (h Hierarchic) Cat(cat string) (map[string]gopium.Struct, bool) {
	// remove root cat from the cat
	cat = strings.Replace(cat, h.rcat, "", 1)
	flat, ok := h.cats[cat]
	return flat, ok
}

// Catflat returns hierarchic categoty
// flat collection if any exists
func (h Hierarchic) Catflat(cat string) (Flat, bool) {
	// cat flat is just alias to cat
	f, ok := h.Cat(cat)
	return Flat(f), ok
}

// Full converts hierarchic collection to flat collection
func (h Hierarchic) Full() map[string]gopium.Struct {
	// collect all structs by key
	flat := make(Flat)
	for _, lsts := range h.cats {
		for key, st := range lsts {
			flat[key] = st
		}
	}
	return flat
}

// Flat converts hierarchic collection to flat collection
func (h Hierarchic) Flat() Flat {
	// flat is just alias to full
	return Flat(h.Full())
}

// Rcat finds root category that
// all other category contain for collection
func (h Hierarchic) Rcat() string {
	// make cats order predictable
	cats := make([]string, 0, len(h.cats))
	for cat := range h.cats {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	// go through cats
	var rcat string
	for _, cat := range cats {
		cat = filepath.Dir(cat)
		switch {
		case cat == ".":
			// just skip empty dir
		case rcat == "":
			// if root cat hasn't been
			// initialized yet set it
			rcat = cat
		case strings.Contains(rcat, cat):
			// if root cat contains cat
			// update root cat
			rcat = cat
		case strings.Contains(cat, rcat):
			// just skip this case
		default:
			// otherwise there are no
			// specific root cat found
			return ""
		}
	}
	return rcat
}

// Len calculates total len of hierarchic collection
func (h Hierarchic) Len() int {
	var l int
	for _, cat := range h.cats {
		l += len(cat)
	}
	return l
}
