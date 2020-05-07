package data

import (
	"context"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/types"
	"path/filepath"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/typepkg"

	"golang.org/x/tools/go/packages"
)

// init types cache map and sync
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

// NewParser creates parser for single tests data
func NewParser(pkg string) gopium.Parser {
	p := typepkg.ParserXToolPackagesAst{
		Path:       fmt.Sprintf("%s/%s", "src/1pkg/gopium/tests/data", pkg),
		Root:       build.Default.GOPATH,
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments | parser.AllErrors,
		BuildFlags: []string{"-tags=tests_data"},
	}
	return Parser{Parser: p}
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
