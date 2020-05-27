package data

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/typepkg"

	"golang.org/x/tools/go/packages"
)

var Gopium string

// sets gopium data path
func init() {
	// grabs running root path
	p, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	// until we rich project root
	for path.Base(p) != gopium.NAME {
		p = path.Dir(p)
	}
	Gopium = p
}

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
		Pattern:    fmt.Sprintf("tests/data/%s", pkg),
		Path:       filepath.Join(Gopium, "tests", "data", pkg),
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

// locator defines tests data locator implementation
// which reuses underlying locator
// but simplifies and purifies ID generation
type locator struct {
	loc gopium.Locator
}

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

// Writer defines tests data writter implementation
// which reuses underlying locator
// but purifies location generation
type Writer struct {
	Writer gopium.CategoryWriter
}

// Generate writer implementation
func (w Writer) Generate(loc string) (io.WriteCloser, error) {
	// purify the loc then
	// just reuse underlying writer
	return w.Writer.Generate(purify(loc))
}

// Category writer implementation
func (w Writer) Category(cat string) error {
	return w.Writer.Category(cat)
}

// purify helps to transform
// absolute path to relative local one
func purify(loc string) string {
	// remove abs part from loc
	// replace os path separators
	// with underscores and trim them
	loc = strings.Replace(loc, Gopium, "", 1)
	loc = strings.ReplaceAll(loc, string(os.PathSeparator), "_")
	return strings.Trim(loc, "_")
}
