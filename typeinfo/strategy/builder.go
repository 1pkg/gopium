package strategy

import (
	"fmt"
	"os"

	"1pkg/gopium"
	"1pkg/gopium/typeinfo"
)

// list of registred strategy names
var (
	TypeInfoOutPrettyJsonStd gopium.StrategyName = "TypeInfoOut-PrettyJsonStd"
)

// Builder defines typeinfo strategy StrategyBuilder implementation
// that uses typeinfo Extractor and related strategies
type Builder struct {
	e typeinfo.Extractor
}

// NewBuilder creates instance of typeinfo strategy Builder
// and requires typeinfo Extractor for it
func NewBuilder(e typeinfo.Extractor) Builder {
	return Builder{e: e}
}

// Build typeinfo strategy StrategyBuilder implementation
func (b Builder) Build(name gopium.StrategyName) (gopium.Strategy, error) {
	switch name {
	case TypeInfoOutPrettyJsonStd:
		return (&typeOut{
			typeMap: typeMap{e: b.e},
			f:       gopium.PrettyJson,
			w:       os.Stdout,
		}).Execute, nil
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", name)
	}
}
