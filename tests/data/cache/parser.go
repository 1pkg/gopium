package cache

import (
	"context"
	"go/ast"
	"go/types"
	"path/filepath"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/typepkg"
)

// prebake types cache map and sync
var (
	tcache map[string]typesloc = make(map[string]typesloc, 5)
	tmutex sync.Mutex
)

// typesloc data transfer object
// that contains types package and loc
type typesloc struct {
	pkg *types.Package
	loc gopium.Locator
}

// parser defines parser implementation
// that adds internal caching for results
type Parser struct {
	Parser gopium.Parser
}

// ParseTypes cache parser implementation
func (p Parser) ParseTypes(ctx context.Context) (*types.Package, gopium.Locator, error) {
	// check that known parser should be cached
	if parser, ok := p.Parser.(typepkg.ParserXToolPackagesAst); ok {
		// access cache syncroniusly
		defer tmutex.Unlock()
		tmutex.Lock()
		// build full dir cache key
		dir := filepath.Join(parser.Root, parser.Path)
		// check if key exists in cache
		if tp, ok := tcache[dir]; ok {
			return tp.pkg, tp.loc, nil
		}
		// if not do actual parsing
		pkg, loc, err := p.Parser.ParseTypes(ctx)
		// store result to cache if no error occured
		if err == nil {
			tcache[dir] = typesloc{pkg: pkg, loc: loc}
		}
		return pkg, loc, err
	}
	// otherwise use real parser
	return p.Parser.ParseTypes(ctx)
}

// ParseAst cache parser implementation
func (p Parser) ParseAst(ctx context.Context) (*ast.Package, gopium.Locator, error) {
	// it's cheap to parse ast each time
	return p.Parser.ParseAst(ctx)
}
