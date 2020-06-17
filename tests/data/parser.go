package data

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"path/filepath"
	"sync"

	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests"
	"github.com/1pkg/gopium/typepkg"

	"golang.org/x/tools/go/packages"
)

// init types cache map and sync
var (
	tcache map[string]typesloc = make(map[string]typesloc, 6)
	tmutex sync.Mutex
)

// typesloc data transfer object
// that contains types package and loc
type typesloc struct {
	loc gopium.Locator `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	pkg *types.Package `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_   [8]byte        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Parser defines tests data parser implementation
// that adds internal caching for results
type Parser struct {
	p gopium.Parser `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// NewParser creates parser for single tests data
func NewParser(pkg string) gopium.Parser {
	p := &typepkg.ParserXToolPackagesAst{
		Pattern:    fmt.Sprintf("github.com/1pkg/gopium/tests/data/%s", pkg),
		Path:       filepath.Join(tests.Gopium, "tests", "data", pkg),
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments | parser.AllErrors,
		BuildFlags: []string{"-tags=tests_data"},
	}
	return Parser{p: p}
}

// ParseTypes parser implementation
func (p Parser) ParseTypes(ctx context.Context, src ...byte) (*types.Package, gopium.Locator, error) {
	// check that known parser should be cached
	if xparser, ok := p.p.(*typepkg.ParserXToolPackagesAst); ok {
		// access cache syncroniusly
		defer tmutex.Unlock()
		tmutex.Lock()
		// build full dir cache key
		dir := filepath.Join(xparser.Root, xparser.Path)
		// check if key exists in cache
		if tp, ok := tcache[dir]; ok {
			return tp.pkg, tp.loc, nil
		}
		// if not then do actual parsing
		// and wrap locator with data locator
		pkg, loc, err := p.p.ParseTypes(ctx, src...)
		// store result to cache if no error occurred
		if err == nil {
			tcache[dir] = typesloc{pkg: pkg, loc: locator{loc: loc}}
		}
		return pkg, locator{loc: loc}, err
	}
	// otherwise use real parser
	// also wrap locator with data locator
	pkg, loc, err := p.p.ParseTypes(ctx, src...)
	return pkg, locator{loc: loc}, err
}

// ParseAst cache parser implementation
func (p Parser) ParseAst(ctx context.Context, src ...byte) (*ast.Package, gopium.Locator, error) {
	// it's cheap to parse ast each time
	// also wrap locator with data locator
	pkg, loc, err := p.p.ParseAst(ctx, src...)
	return pkg, locator{loc: loc}, err
}
