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

// init global cache maps
var gtypes, gast *sync.Map = &sync.Map{}, &sync.Map{}

// typesloc data transfer object
// that contains types package and loc
type typesloc struct {
	pkg *types.Package
	loc gopium.Locator
}

// astloc data transfer object
// that contains ast package and loc
type astloc struct {
	pkg *ast.Package
	loc gopium.Locator
}

// parser defines parser implementation
// that adds internal caching for results
type Parser struct {
	parser gopium.Parser
	types  *sync.Map
	ast    *sync.Map
}

// With creates new instance of cached parser
// with shared cache and provided parser
func (p Parser) With(parser gopium.Parser) gopium.Parser {
	// create new cached parser
	return &Parser{
		parser: parser,
		types:  gtypes,
		ast:    gast,
	}
}

// ParseTypes cache parser implementation
func (p Parser) ParseTypes(ctx context.Context) (*types.Package, gopium.Locator, error) {
	// check that known parser should be cached
	if parser, ok := p.parser.(typepkg.ParserXToolPackagesAst); ok {
		// build full dir cache key
		dir := filepath.Join(parser.Root, parser.Path)
		// check if key exists in cache
		if val, ok := p.types.Load(dir); ok {
			if pl, ok := val.(typesloc); ok {
				return pl.pkg, pl.loc, nil
			}
		}
		// if not do actual parsing
		pkg, loc, err := p.parser.ParseTypes(ctx)
		// store result to cache if no error occured
		if err == nil {
			p.types.Store(dir, typesloc{pkg: pkg, loc: loc})
		}
		return pkg, loc, err
	}
	return p.parser.ParseTypes(ctx)
}

// ParseAst cache parser implementation
func (p Parser) ParseAst(ctx context.Context) (*ast.Package, gopium.Locator, error) {
	// check that known parser should be cached
	if parser, ok := p.parser.(typepkg.ParserXToolPackagesAst); ok {
		// build full dir cache key
		dir := filepath.Join(parser.Root, parser.Path)
		// check if key exists in cache
		if val, ok := p.ast.Load(dir); ok {
			if al, ok := val.(astloc); ok {
				return al.pkg, al.loc, nil
			}
		}
		// if not do actual parsing
		pkg, loc, err := p.parser.ParseAst(ctx)
		// store result to cache if no error occured
		if err == nil {
			p.ast.Store(dir, astloc{pkg: pkg, loc: loc})
		}
		return pkg, loc, err
	}
	return p.parser.ParseAst(ctx)
}
