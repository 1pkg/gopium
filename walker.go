package gopium

import (
	"context"
	"regexp"
)

// Walker defines hierarchical walker abstraction
// that applies some strategy to tree structures
type Walker interface {
	VisitTop(context.Context, *regexp.Regexp, Strategy)  // only top level of the tree
	VisitDeep(context.Context, *regexp.Regexp, Strategy) // all levels of the tree
}
