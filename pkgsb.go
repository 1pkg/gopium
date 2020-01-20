package main

import (
	"fmt"
)

// SgBuilder defines builder abstraction
// that helps to create strategy by name
type SgBuilder interface {
	Build(string) (Strategy, error)
}

// Pkgsb defines package strategy builder implementation
// that uses type info extractor abstraction to build strategies
type Pkgsb TiExt

// Build package strategy builder implementation
func (sb Pkgsb) Build(sg string) (Strategy, error) {
	switch sg {
	default:
		return nil, fmt.Errorf("strategy `%s` wasn't found", sg)
	}
}
