package gopium

import (
	"context"
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
