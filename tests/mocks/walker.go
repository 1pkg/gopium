package mocks

import (
	"context"
	"regexp"

	"1pkg/gopium"
)

// Walker defines mock walker implementation
type Walker struct {
	Err error
}

// Visit mock implementation
func (w Walker) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	// check error at start
	if w.Err != nil {
		return w.Err
	}
	// return context error otherwise
	return ctx.Err()
}

// WalkerBuilder defines mock walker builder implementation
type WalkerBuilder struct {
	Walker gopium.Walker
	Err    error
}

// Build mock implementation
func (b WalkerBuilder) Build(gopium.WalkerName) (gopium.Walker, error) {
	return b.Walker, b.Err
}
