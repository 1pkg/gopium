package strategy

import (
	"context"
	"sort"

	"1pkg/gopium"
)

// lex defines strategy implementation
// that sorts fields accordingly to their names
// in ascending order
type lex struct{}

// Apply lex implementation
func (stg lex) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// then execute lexicographical sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		return r.Fields[i].Name < r.Fields[j].Name
	})
	return
}
