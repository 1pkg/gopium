package strategy

import (
	"context"
	"go/types"
	"sort"

	"1pkg/gopium"
)

// lexicographical defines struct fields names
// lexicographical sorting Strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then sorts fields accordingly to their names in ascending order
type lexicographical struct {
	m gopium.Maven
}

// Apply lexicographical implementation
func (stg lexicographical) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// then execute lexicographical sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		return r.Fields[i].Name < r.Fields[j].Name
	})
	return
}
