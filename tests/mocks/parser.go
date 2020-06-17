package mocks

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/1pkg/gopium/gopium"
)

// Pos defines mock pos
// data transfer object
type Pos struct {
	ID  string `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Loc string `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Locator defines mock locator implementation
type Locator struct {
	Poses map[token.Pos]Pos `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 8 bytes; struct align: 8 bytes; struct aligned size: 8 bytes; - 🌺 gopium @1pkg

// ID mock implementation
func (l Locator) ID(pos token.Pos) string {
	// check if we have it in vals
	if t, ok := l.Poses[pos]; ok {
		return t.ID
	}
	// otherwise return default val
	return ""
}

// Loc mock implementation
func (l Locator) Loc(pos token.Pos) string {
	// check if we have it in vals
	if t, ok := l.Poses[pos]; ok {
		return t.Loc
	}
	// otherwise return default val
	return ""
}

// Locator mock implementation
func (l Locator) Locator(string) (gopium.Locator, bool) {
	return l, true
}

// Fset mock implementation
func (l Locator) Fset(string, *token.FileSet) (*token.FileSet, bool) {
	return token.NewFileSet(), true
}

// Root mock implementation
func (l Locator) Root() *token.FileSet {
	return token.NewFileSet()
}

// Parser defines mock parser implementation
type Parser struct {
	Parser   gopium.Parser `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Typeserr error         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Asterr   error         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_        [16]byte      `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// ParseTypes mock implementation
func (p Parser) ParseTypes(ctx context.Context, src ...byte) (*types.Package, gopium.Locator, error) {
	// if parser provided use it
	if p.Parser != nil {
		pkg, loc, _ := p.Parser.ParseTypes(ctx, src...)
		return pkg, loc, p.Typeserr
	}
	return types.NewPackage("", ""), Locator{}, p.Typeserr
}

// ParseAst mock implementation
func (p Parser) ParseAst(ctx context.Context, src ...byte) (*ast.Package, gopium.Locator, error) {
	// if parser provided use it
	if p.Parser != nil {
		pkg, loc, _ := p.Parser.ParseAst(ctx, src...)
		return pkg, loc, p.Asterr
	}
	return &ast.Package{}, Locator{}, p.Asterr
}
