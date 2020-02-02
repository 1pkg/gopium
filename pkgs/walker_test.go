package pkgs

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"regexp"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/fmts"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

// LoadWalkerTestCase helps to load
// pkgs tests data directly from source code package
// using only test case name
func LoadWalkerTestCase(tcase string) (gopium.Walker, error) {
	regex, _ := regexp.Compile(`.*`)
	p := ParserXTool{
		Patterns: []string{"1pkg/gopium/pkgs/pkgs_test_data/" + tcase},
		LoadMode: packages.LoadAllSyntax,
	}.Parse
	return NewWalker(context.Background(), fmts.FullName, regex, p)
}

func TestNewWalker(t *testing.T) {
	regex, _ := regexp.Compile(`^foobar$`)
	table := []struct {
		name  string
		ctx   context.Context
		hn    fmts.HierarchyName
		regex *regexp.Regexp
		p     Parser
		w     Walker
		err   error
	}{
		{
			name:  "nil fmts.HierarchyName, NewWalker should return error",
			ctx:   context.Background(),
			hn:    nil,
			regex: regex,
			p: ParserMock{
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
				fset: token.NewFileSet(),
			}.Parse,
			w:   Walker{},
			err: errors.New("hierarchy name wasn't defined"),
		},
		{
			name:  "packages Parser returns error, NewWalker should just pass it",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: regex,
			p:     ParserError("error package test error").Parse,
			w:     Walker{},
			err:   errors.New("error package test error"),
		},
		{
			name:  "package name wasn't found, NewWalker should return error",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: regex,
			p:     ParserNil{}.Parse,
			w:     Walker{},
			err:   errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name:  "packages Parser returns nil type package, NewWalker should return error",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: regex,
			p:     ParserMock{fset: token.NewFileSet()}.Parse,
			w:     Walker{},
			err:   errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name:  "packages Parser returns nil fset, NewWalker should return error",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: regex,
			p:     ParserMock{pkgs: []*types.Package{types.NewPackage("/", "foobar")}}.Parse,
			w:     Walker{},
			err:   errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name:  "package was found, NewWalker should return correct package walker",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: regex,
			p: ParserMock{
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
				fset: token.NewFileSet(),
			}.Parse,
			w: Walker{
				hn:   fmts.FlatName,
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
				fset: token.NewFileSet(),
			},
			err: nil,
		},
	}
	t.Run("NewWalker should return correct results for all cases", func(t *testing.T) {
		for _, tcase := range table {
			t.Run(tcase.name, func(t *testing.T) {
				w, err := NewWalker(tcase.ctx, tcase.hn, tcase.regex, tcase.p)
				// "github.com/stretchr/testify/assert" can't be used for function equality.
				// `Function equality cannot be determined and will always fail.`
				w.hn, tcase.w.hn = nil, nil
				assert.Equal(t, tcase.w, w)
				assert.Equal(t, tcase.err, err)
			})
		}
	})
}

func TestWalkerVisitTopErrorExplode(t *testing.T) {
	regex, _ := regexp.Compile(`^Foo`)
	w, err := LoadWalkerTestCase("walker_visit_top_error_explode")
	assert.NoError(t, err)
	stg := gopium.StrategyError("this should never happen")
	assert.Panics(t, func() {
		w.VisitTop(context.Background(), regex, stg.Execute)
	})
}

func TestWalkerVisitDeepErrorExplode(t *testing.T) {
	regex, _ := regexp.Compile(`Bar$`)
	w, err := LoadWalkerTestCase("walker_visit_deep_error_explode")
	assert.NoError(t, err)
	stg := gopium.StrategyError("this should never happen")
	assert.Panics(t, func() {
		w.VisitTop(context.Background(), regex, stg.Execute)
	})
}
