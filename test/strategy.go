// +build test

package test

import (
	"context"
	"go/types"

	"1pkg/gopium"
)

// StrategyMock defines strategy mock implementation
type StrategyMock struct{}

// Apply StrategyMock implementation
func (stg StrategyMock) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// set structure name
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get tag
		tag := st.Tag(i)
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     f.Type().String(),
			Tag:      tag,
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	o = r
	return
}

// StrategyError defines strategy error implementation
type StrategyError struct {
	err error
}

// Apply StrategyError implementation
func (stg StrategyError) Apply(ctx context.Context, name string, st *types.Struct) (o gopium.Struct, r gopium.Struct, err error) {
	// just set error
	err = stg.err
	return
}
