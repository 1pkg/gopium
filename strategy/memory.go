package strategy

import (
	"context"
	"go/types"
	"sort"

	"1pkg/gopium"
)

// memory defines struct optimal memory fields sorting strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then sorts fields accordingly to their sizes in descending order
type memory struct {
	m gopium.Maven
}

// Apply memory implementation
func (stg memory) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// then execute memory sorting
	sort.SliceStable(r.Fields, func(i, j int) bool {
		return r.Fields[j].Size < r.Fields[i].Size
	})
	return
}
