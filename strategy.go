package gopium

import (
	"context"
	"fmt"
	"go/types"
)

// Strategy defines "pure func" action abstraction
// that applies some strategy payload on types.Struct an it's name
// and returns resulted Struct object or error
type Strategy interface {
	Apply(ctx context.Context, name string, st *types.Struct) StructError
}

// StrategyName defines registred strategy name abstraction
// used by StrategyBuilder to build registred strategies
type StrategyName string

// StrategyBuilder defines strategy builder abstraction
// that helps to create Strategy by StrategyName
type StrategyBuilder interface {
	Build(StrategyName) (Strategy, error)
}

// StrategyMock defines Strategy mock implementation
type StrategyMock struct{}

// Apply StrategyMock implementation
func (stg StrategyMock) Apply(ctx context.Context, name string, st *types.Struct) (r StructError) {
	// build full hierarchical name of the structure
	r.Struct.Name = fmt.Sprintf("%s/%s", name, st)
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Struct.Fields = make([]Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get tag
		tag := st.Tag(i)
		// fill field structure
		r.Struct.Fields = append(r.Struct.Fields, Field{
			Name:     f.Name(),
			Type:     f.Type().String(),
			Size:     0,
			Tag:      tag,
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return
}

// StrategyError defines Strategy error implementation
type StrategyError struct {
	err error
}

// Apply StrategyError implementation
func (stg StrategyError) Apply(ctx context.Context, name string, st *types.Struct) (r StructError) {
	// just set error
	r.Error = stg.err
	return
}
