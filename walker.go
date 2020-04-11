package gopium

import (
	"context"
	"regexp"
)

// Walker defines hierarchical walker abstraction
// that applies some strategy to tree structures
// and modifies tree or creates other side effects
type Walker interface {
	Visit(context.Context, *regexp.Regexp, Strategy, bool) error
}

// WalkerName defines registered walker name abstraction
// used by walker builder to build registered walkers
type WalkerName string

// WalkerBuilder defines walker builder abstraction
// that helps to create walker by walker name
type WalkerBuilder interface {
	Build(WalkerName) (Walker, error)
}
