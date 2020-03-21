package strategy

import (
	"context"
	"fmt"
	"go/types"

	"1pkg/gopium"
)

// separating_caching defines system struct padding fields strategy implementation
// that uses enum strategy to get gopium.Field DTO for each field
// then adds cache padding field before and cache padding field after structure fields list
type separating_caching struct {
	m gopium.Maven
	l uint
}

// Apply separating_caching implementation
func (stg separating_caching) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// first apply enum strategy
	enum := enum{stg.m}
	o, r, err = enum.Apply(ctx, name, st)
	// add field before and after
	// with sys align size
	sep := stg.m.SysCache(stg.l)
	r.Fields = append([]gopium.Field{
		gopium.Field{
			Name:  "_",
			Type:  fmt.Sprintf("[%d]byte", sep),
			Size:  sep,
			Align: 1,
		},
	}, r.Fields...)
	r.Fields = append(r.Fields, gopium.Field{
		Name:  "_",
		Type:  fmt.Sprintf("[%d]byte", sep),
		Size:  sep,
		Align: 1,
	})
	return
}
