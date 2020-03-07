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

// StrategyMode defines registered strategy mode abstraction
// used by StrategyBuilder to build registered strategies
type StrategyMode uint

// Has checks if mode includes sub mode
func (mode StrategyMode) Has(m StrategyMode) bool {
	return mode&m == m
}

// StrategyBuilder defines strategy builder abstraction
// that helps to create strategy by strategy name
type StrategyBuilder interface {
	Build(StrategyName, StrategyMode) (Strategy, error)
}
