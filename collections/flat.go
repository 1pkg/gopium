package collections

import (
	"fmt"
	"sort"
	"strings"

	"1pkg/gopium"
)

// Flat defines strucs flat collection
// which is categorized only by id
type Flat map[string]gopium.Struct

// Sorted converts flat collection
// to sorted slice of structs
// note: it's possible due to next:
// generated ids are ordered inside same loc
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
		var numi, numj int
		var sumi, sumj string
		// scanf works only with space
		// separated values so we need
		// to apply this format first
		//
		// in case of any pattern error
		// just apply natural sort
		// otherwise sort it by id
		_, erri := fmt.Sscanf(strings.Replace(ids[i], ":", " ", 1), "%s %d", &sumi, &numi)
		_, errj := fmt.Sscanf(strings.Replace(ids[j], ":", " ", 1), "%s %d", &sumj, &numj)
		switch {
		case erri != nil && errj != nil:
			return ids[i] < ids[j]
		case erri != nil:
			return false
		case errj != nil:
			return true
		default:
			return numi < numj
		}
	})
	// collect all structs in asc order
	for _, id := range ids {
		sorted = append(sorted, f[id])
	}
	return sorted
}
