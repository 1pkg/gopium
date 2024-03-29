package data

import (
	"fmt"
	"go/token"

	"github.com/1pkg/gopium/gopium"
)

// locator defines tests data locator implementation
// which reuses underlying locator
// but simplifies and purifies ID generation
type locator struct {
	loc gopium.Locator `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// ID locator implementation
func (l locator) ID(p token.Pos) string {
	// check if such file exists
	if f := l.loc.Root().File(p); f != nil {
		// purify the loc then
		// generate ordered id
		return fmt.Sprintf("%s:%d", purify(f.Name()), f.Line(p))
	}
	return ""
}

// Loc locator implementation
func (l locator) Loc(p token.Pos) string {
	return l.loc.Loc(p)
}

// Locator locator implementation
func (l locator) Locator(loc string) (gopium.Locator, bool) {
	return l.loc.Locator(loc)
}

// Fset locator implementation
func (l locator) Fset(loc string, fset *token.FileSet) (*token.FileSet, bool) {
	return l.loc.Fset(loc, fset)
}

// Root locator implementation
func (l locator) Root() *token.FileSet {
	return l.loc.Root()
}
