package pkgs

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"path/filepath"
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
func LoadWalkerTestCase(tcases ...string) (Walker, error) {
	// setup abs path and regex
	abs, _ := filepath.Abs("./")
	regex, _ := regexp.Compile(`.*`)
	// patch all test cases names
	for i, tcase := range tcases {
		tcases[i] = "1pkg/gopium/pkgs/pkgs_test_data/" + tcase
	}
	// setup ParserXTool configs
	p := ParserXTool{
		Patterns:   tcases,
		LoadMode:   packages.LoadAllSyntax,
		BuildFlags: []string{"-tags=pkgs_test_data"},
	}.Parse
	return NewWalker(context.Background(), fmts.Root(abs).FullName, regex, p)
}

func TestNewWalker(t *testing.T) {
	table := []struct {
		name  string
		ctx   context.Context
		hn    fmts.HierarchyName
		regex string
		p     Parser
		w     Walker
		err   error
	}{
		{
			name:  "nil fmts.HierarchyName, NewWalker should return error",
			ctx:   context.Background(),
			hn:    nil,
			regex: `^foobar$`,
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
			regex: `^foobar$`,
			p:     ParserError("error package test error").Parse,
			w:     Walker{},
			err:   errors.New("error package test error"),
		},
		{
			name:  "package name wasn't found, NewWalker should return error",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: `^foobar$`,
			p:     ParserNil{}.Parse,
			w:     Walker{},
			err:   errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name:  "packages Parser returns nil type package, NewWalker should return error",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: `^foobar$`,
			p:     ParserMock{fset: token.NewFileSet()}.Parse,
			w:     Walker{},
			err:   errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name:  "packages Parser returns nil fset, NewWalker should return error",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: `^foobar$`,
			p:     ParserMock{pkgs: []*types.Package{types.NewPackage("/", "foobar")}}.Parse,
			w:     Walker{},
			err:   errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name:  "package was found, NewWalker should return correct package walker",
			ctx:   context.Background(),
			hn:    fmts.FlatName,
			regex: `^foobar$`,
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
				// prepare regex
				regex, err := regexp.Compile(tcase.regex)
				assert.NoError(t, err)
				// execute NewWalker action
				w, err := NewWalker(tcase.ctx, tcase.hn, regex, tcase.p)
				// patch Walker fmts.HierarchyName func as
				// "github.com/stretchr/testify/assert" can't be used for function equality.
				// `Function equality cannot be determined and will always fail.`
				w.hn, tcase.w.hn = nil, nil
				// check result
				assert.Equal(t, tcase.w, w)
				assert.Equal(t, tcase.err, err)
			})
		}
	})
}

func TestWalkerVisitTopErrorExplode(t *testing.T) {
	// prepare error theoretically ErrorStrategy and regex
	regex, err := regexp.Compile(`^Foo`)
	assert.NoError(t, err)
	w, err := LoadWalkerTestCase("walker_visit_top_error_explode")
	assert.NoError(t, err)
	stg := gopium.StrategyError("this should never happen")
	// execute action and check panic outcome
	assert.Panics(t, func() {
		w.VisitTop(context.Background(), regex, stg.Execute)
	})
}

func TestWalkerVisitDeepErrorExplode(t *testing.T) {
	// prepare error theoretically ErrorStrategy and regex
	regex, err := regexp.Compile(`Bar$`)
	assert.NoError(t, err)
	w, err := LoadWalkerTestCase("walker_visit_deep_error_explode")
	assert.NoError(t, err)
	stg := gopium.StrategyError("this should never happen")
	// execute action and check panic outcome
	assert.Panics(t, func() {
		w.VisitTop(context.Background(), regex, stg.Execute)
	})
}

func TestWalkerVisitTop(t *testing.T) {
	table := []struct {
		name   string
		ctx    context.Context
		lwtc   []string
		regex  string
		result map[string]string
	}{
		{
			name:   "empty packages list, Walker VisitTop should process nothing",
			ctx:    context.Background(),
			lwtc:   []string{},
			regex:  `.*`,
			result: make(map[string]string),
		},
		{
			name:   "empty structs list, Walker VisitTop should process nothing",
			ctx:    context.Background(),
			lwtc:   []string{"test_walker_visit_top_empty"},
			regex:  `.*`,
			result: make(map[string]string),
		},
		{
			name: "empty structs list even in multiple packages, Walker VisitTop should process nothing",
			ctx:  context.Background(),
			lwtc: []string{
				"test_walker_visit_top_multiple_empty/p1",
				"test_walker_visit_top_multiple_empty/p2",
				"test_walker_visit_top_multiple_empty/p3",
			},
			regex:  `.*`,
			result: make(map[string]string),
		},
		{
			name:   "nested structs only list, Walker VisitTop should process nothing",
			ctx:    context.Background(),
			lwtc:   []string{"test_walker_visit_top_nested_structs_only"},
			regex:  `.*`,
			result: make(map[string]string),
		},
		{
			name:  "single top level structs list, Walker VisitTop should process it",
			ctx:   context.Background(),
			lwtc:  []string{"test_walker_visit_top_single_top_level_struct"},
			regex: `.*`,
			result: map[string]string{
				"./pkgs_test_data/test_walker_visit_top_single_top_level_struct/test.go:L6 FooBar": "struct{xint int; xstring string}",
			},
		},
		{
			name:  "multiple top level structs list, Walker VisitTop should filter and process them",
			ctx:   context.Background(),
			lwtc:  []string{"test_walker_visit_top_multiple_top_level_structs"},
			regex: `^Foo.*`,
			result: map[string]string{
				"./pkgs_test_data/test_walker_visit_top_multiple_top_level_structs/test_1.go:L6 FooBar": "struct{xint int; xstring string}",
				"./pkgs_test_data/test_walker_visit_top_multiple_top_level_structs/test_2.go:L9 Fooooo": "struct{teststring string}",
			},
		},
		{
			name:  "mixed structs list, Walker VisitTop should filter and process only top level structs",
			ctx:   context.Background(),
			lwtc:  []string{"test_walker_visit_top_mixed_level_structs"},
			regex: `^FooBar.*`,
			result: map[string]string{
				"./pkgs_test_data/test_walker_visit_top_mixed_level_structs/test.go:L36 FooBarDouble": "struct{f 1pkg/gopium/pkgs/pkgs_test_data/test_walker_visit_top_mixed_level_structs.FooBar; s 1pkg/gopium/pkgs/pkgs_test_data/test_walker_visit_top_mixed_level_structs.FooBar}",
				"./pkgs_test_data/test_walker_visit_top_mixed_level_structs/test.go:L6 FooBar":        "struct{xint int; xstring string}",
			},
		},
		{
			name:  "mixed declarations, Walker VisitTop should filter and process only top level structs",
			ctx:   context.Background(),
			lwtc:  []string{"test_walker_visit_top_mixed_declarations"},
			regex: `^FooBar.*`,
			result: map[string]string{
				"./pkgs_test_data/test_walker_visit_top_mixed_declarations/test.go:L6 FooBarStruct":  "struct{structstring string}",
				"./pkgs_test_data/test_walker_visit_top_mixed_declarations/test.go:L10 FooBarCoType": "struct{structstring string}", // TODO shouldn't be processed
			},
		},
	}
	t.Run("Walker VisitTop should process correctly all cases", func(t *testing.T) {
		for _, tcase := range table {
			t.Run(tcase.name, func(t *testing.T) {
				// create walker
				var w Walker
				if len(tcase.lwtc) > 0 {
					// if load walker test cases list isn't empty
					var err error
					w, err = LoadWalkerTestCase(tcase.lwtc...)
					assert.NoError(t, err)
				} else {
					// otherwise create empty walker
					w = Walker{
						hn:   fmts.FlatName,
						pkgs: []*types.Package{types.NewPackage("/", "foobar")},
						fset: token.NewFileSet(),
					}
				}
				// prepare Strategy and regex
				regex, err := regexp.Compile(tcase.regex)
				assert.NoError(t, err)
				stg := make(gopium.StrategyMock)
				// execute action and check result
				w.VisitTop(tcase.ctx, regex, stg.Execute)
				assert.Equal(t, tcase.result, map[string]string(stg))
			})
		}
	})
}
