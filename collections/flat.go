package collections

import (
	"fmt"
	"sort"

	"1pkg/gopium"
)

// Flat defines strucs flat collection
// which is categorized only by id
type Flat map[string]gopium.Struct

// Sorted converts flat collection
// to sorted slice of structs
// note: it's possible due to next:
// generated id would be ordered inside same loc
func (f Flat) Sorted() []gopium.Struct {
	// preapare ids and sorted slice
	ids := make([]string, 0, len(f))
	sorted := make([]gopium.Struct, 0, len(f))
	// collect all ids
	for id := range f {
		ids = append(ids, id)
	}
	// sort all ids in asc order
	sort.SliceStable(ids, func(i, j int) bool {
		// only first part of "%d-%s" id is ordered
		// so we need to parse and compare it
		var idi, idj int
		var sumi, sumj string
		// in case of any pattern error
		// just apply natural sort
		_, erri := fmt.Sscanf(ids[i], "%d-%s", &idi, &sumi)
		_, errj := fmt.Sscanf(ids[j], "%d-%s", &idj, &sumj)
		if erri != nil && errj != nil {
			return ids[i] < ids[j]
		} else if erri != nil {
			return false
		} else if errj != nil {
			return true
		} else {
			return idi < idj
		}
	})
	// collect all structs in asc order
	for _, id := range ids {
		sorted = append(sorted, f[id])
	}
	return sorted
}
