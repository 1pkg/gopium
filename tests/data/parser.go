package data

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"path/filepath"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/tests"
	"1pkg/gopium/typepkg"

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
	pkg *types.Package
	loc gopium.Locator
}

// Parser defines tests data parser implementation
// that adds internal caching for results
type Parser struct {
	p gopium.Parser
}

// NewParser creates parser for single tests data
func NewParser(pkg string) gopium.Parser {
	p := &typepkg.ParserXToolPackagesAst{
		Pattern:    fmt.Sprintf("1pkg/gopium/tests/data/%s", pkg),
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
