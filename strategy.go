package gopium

import (
	"context"
	"go/types"
)

// Strategy defines "pure func" action abstraction
// that applies some strategy payload on types.Struct an it's name
// and returns resulted Struct object or error
type Strategy interface {
	Apply(ctx context.Context, name string, st *types.Struct) (o Struct, r Struct, err error)
}

// StrategyName defines registred Strategy name abstraction
// used by StrategyBuilder to build registred strategies
type StrategyName string

// StrategyBuilder defines Strategy builder abstraction
// that helps to create Strategy by StrategyName
type StrategyBuilder interface {
	Build(StrategyName) (Strategy, error)
}

// StrategyMock defines Strategy mock implementation
type StrategyMock struct{}

// Apply StrategyMock implementation
func (stg StrategyMock) Apply(ctx context.Context, name string, st *types.Struct) (o Struct, r Struct, err error) {
	// build full hierarchical name of the structure
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

// StrategyError defines Strategy error implementation
type StrategyError struct {
	err error
}

// Apply StrategyError implementation
func (stg StrategyError) Apply(ctx context.Context, name string, st *types.Struct) (o Struct, r Struct, err error) {
	// just set error
	err = stg.err
	return
}
