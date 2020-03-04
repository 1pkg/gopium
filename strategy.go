package gopium

import (
	"context"
	"go/types"
)

// Strategy defines action abstraction
// that applies some strategy payload on types.Struct
// and returns origin, resulted struct objects or error
type Strategy interface {
	Apply(ctx context.Context, name string, st *types.Struct) (o Struct, r Struct, err error)
}

// StrategyName defines registered strategy name abstraction
// used by StrategyBuilder to build registered strategies
type StrategyName string

// StrategyBuilder defines strategy builder abstraction
// that helps to create strategy by strategy name
type StrategyBuilder interface {
	Build(StrategyName) (Strategy, error)
}

// StrategyMock defines strategy mock implementation
type StrategyMock struct{}

// Apply StrategyMock implementation
func (stg StrategyMock) Apply(ctx context.Context, name string, st *types.Struct) (o Struct, r Struct, err error) {
	// set structure name
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get tag
		tag := st.Tag(i)
		// fill field structure
		r.Fields = append(r.Fields, Field{
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
func (stg StrategyError) Apply(ctx context.Context, name string, st *types.Struct) (o Struct, r Struct, err error) {
	// just set error
	err = stg.err
	return
}
