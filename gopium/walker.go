package gopium

import (
	"context"
	"regexp"
)

// Walker defines hierarchical walker abstraction
// that applies some strategy to code tree structures
// and modifies them or creates other related side effects
type Walker interface {
	Visit(context.Context, *regexp.Regexp, Strategy) error
}

// WalkerName defines registered walker name abstraction
// used by walker builder to build registered walkers
type WalkerName string

// WalkerBuilder defines walker builder abstraction
// that helps to create single walker by walker name
type WalkerBuilder interface {
	Build(WalkerName) (Walker, error)
}
