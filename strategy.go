package gopium

import "context"

// Strategy defines action abstraction
// that applies some strategy payload on types.Struct
// and returns origin, resulted struct objects or error
type Strategy interface {
	Apply(context.Context, Struct) (Struct, error)
}

// StrategyName defines registered strategy name abstraction
// used by StrategyBuilder to build registered strategies
type StrategyName string

// StrategyBuilder defines strategy builder abstraction
// that helps to create strategy by strategy name
type StrategyBuilder interface {
	Build(StrategyName) (Strategy, error)
}
