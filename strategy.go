package gopium

import "context"

// Strategy defines custom action abstraction
// that applies some action payload on struct
// and returns resulted struct object or error
type Strategy interface {
	Apply(context.Context, Struct) (Struct, error)
}

// StrategyName defines registered strategy name abstraction
// used by StrategyBuilder to build registered strategies
type StrategyName string

// StrategyBuilder defines strategy builder abstraction
// that helps to create single strategy by strategies names
type StrategyBuilder interface {
	Build(...StrategyName) (Strategy, error)
}
