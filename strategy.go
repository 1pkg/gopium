package gopium

import (
	"context"
	"errors"
	"go/token"
	"go/types"
)

// Strategy defines action abstraction
// that applies some strategy on types.Struct
type Strategy func(context.Context, string, *types.Struct, *token.FileSet) error

// StrategyName defines known strategy name type
type StrategyName string

// StrategyBuilder defines strategy builder abstraction
// that helps to create Strategy by name
type StrategyBuilder interface {
	Build(StrategyName) (Strategy, error)
}

// StrategyMock defines Strategy mock implementation
type StrategyMock map[string]string

// Execute Strategy mock implementation
func (stg StrategyMock) Execute(ctx context.Context, nm string, st *types.Struct, fset *token.FileSet) error {
	stg[nm] = st.String()
	return nil
}

// StrategyMock defines Strategy error implementation
type StrategyError string

// Execute Strategy error implementation
func (stg StrategyError) Execute(ctx context.Context, nm string, st *types.Struct, fset *token.FileSet) error {
	return errors.New(string(stg))
}
