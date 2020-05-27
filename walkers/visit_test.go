package walkers

import (
	"context"
	"go/token"
	"go/types"
	"reflect"
	"regexp"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/strategies"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
	"1pkg/gopium/typepkg"
)

func TestWithVisit(t *testing.T) {
	// prepare
	table := map[string]struct {
		exp  gopium.Exposer
		loc  gopium.Locator
		bref bool
		r    *regexp.Regexp
		stg  gopium.Strategy
		ch   appliedCh
		deep bool
		ctx  context.Context
		s    *types.Scope
	}{
		"with visit should return expected govisit func": {
			exp:  mocks.Maven{},
			loc:  mocks.Locator{},
			bref: true,
			r:    regexp.MustCompile(`.*`),
			stg:  &mocks.Strategy{},
			ch:   make(appliedCh),
			deep: true,
			ctx:  context.Background(),
			s:    &types.Scope{},
		},
		"with visit should return expected govisit func without bref flag": {
			exp:  mocks.Maven{},
			loc:  mocks.Locator{},
			bref: false,
			r:    regexp.MustCompile(`.*`),
			stg:  &mocks.Strategy{},
			ch:   make(appliedCh),
			deep: true,
			ctx:  context.Background(),
			s:    &types.Scope{},
		},
		"with visit should return expected govisit func without deep flag": {
			exp:  mocks.Maven{},
			loc:  mocks.Locator{},
			bref: true,
			r:    regexp.MustCompile(`.*`),
			stg:  &mocks.Strategy{},
			ch:   make(appliedCh),
			deep: false,
			ctx:  context.Background(),
			s:    &types.Scope{},
		},
		"with visit should return expected govisit func without all flags": {
			exp: mocks.Maven{},
			loc: mocks.Locator{},
			r:   regexp.MustCompile(`.*`),
			stg: &mocks.Strategy{},
			ch:  make(appliedCh),
			ctx: context.Background(),
			s:   &types.Scope{},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			gvisit := with(tcase.exp, tcase.loc, tcase.bref).
				visit(tcase.r, tcase.stg, tcase.ch, tcase.deep)
			gvisit(tcase.ctx, tcase.s)
			// check
			// we can't compare functions directly in go
			// so apply this hack to compare with nil
			if reflect.DeepEqual(gvisit, nil) {
				t.Errorf("actual %v doesn't equal to expected not %v", gvisit, nil)
			}
		})
	}
}

func TestVscope(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	b := strategies.Builder{}
	np, err := b.Build(strategies.Nope)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pck, err := b.Build(strategies.Pack)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 64, 64, 64)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		ctx context.Context
		r   *regexp.Regexp
		m   gopium.Maven
		p   gopium.TypeParser
		loc gopium.Locator
		stg gopium.Strategy
		sts map[string]gopium.Struct
		err error
	}{
		"empty pkg should visit nothing": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("empty"),
			stg: np,
			sts: make(map[string]gopium.Struct),
		},
		"single struct pkg should visit the struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_single_file.go:5": gopium.Struct{
					Name: "Single",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "string",
							Size:     16,
							Align:    8,
							Exported: true,
						},
						{
							Name:     "B",
							Type:     "string",
							Size:     16,
							Align:    8,
							Exported: true,
						},
						{
							Name:     "C",
							Type:     "string",
							Size:     16,
							Align:    8,
							Exported: true,
						},
					},
				},
			},
		},
		"single struct pkg should visit nothing on canceled context": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
			stg: np,
			sts: make(map[string]gopium.Struct),
			err: context.Canceled,
		},
		"single struct pkg should visit nothing on canceled context in closures": {
			ctx: &mocks.Context{After: 2},
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
			stg: np,
			sts: make(map[string]gopium.Struct),
			err: context.Canceled,
		},
		"flat struct pkg should visit all structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_flat_file.go:10": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:16": gopium.Struct{
					Name: "b",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/flat.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:21": gopium.Struct{
					Name: "C",
					Fields: []gopium.Field{
						{
							Name:  "c",
							Type:  "[]string",
							Size:  24,
							Align: 8,
						},
						{
							Name:     "A",
							Type:     "struct{b 1pkg/gopium/tests/data/flat.b; z 1pkg/gopium/tests/data/flat.A}",
							Size:     24,
							Align:    8,
							Exported: true,
						},
					},
				},
				"tests_data_flat_file.go:29": gopium.Struct{
					Name: "c1",
					Fields: []gopium.Field{
						{
							Name:  "c",
							Type:  "[]string",
							Size:  24,
							Align: 8,
						},
						{
							Name:     "A",
							Type:     "struct{b 1pkg/gopium/tests/data/flat.b; z 1pkg/gopium/tests/data/flat.A}",
							Size:     24,
							Align:    8,
							Exported: true,
						},
					},
				},
				"tests_data_flat_file.go:32": gopium.Struct{
					Name: "D",
					Fields: []gopium.Field{
						{
							Name:  "t",
							Type:  "[13]byte",
							Size:  13,
							Align: 1,
						},
						{
							Name:  "b",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:  "_",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:41": gopium.Struct{
					Name: "AZ",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/flat.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
		"flat struct pkg should visit the struct on same loc": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
			loc: mocks.Locator{
				Poses: map[token.Pos]mocks.Pos{
					token.Pos(1799681): {ID: "test-1", Loc: "test"},
					token.Pos(1799769): {ID: "test-2", Loc: "test"},
					token.Pos(1799802): {ID: "test-1", Loc: "test"},
					token.Pos(1799860): {ID: "test-2", Loc: "test"},
					token.Pos(1799915): {ID: "test-1", Loc: "test"},
					token.Pos(1800016): {ID: "test-1", Loc: "test"},

					token.Pos(1953758): {ID: "test-1", Loc: "test"},
					token.Pos(1953846): {ID: "test-2", Loc: "test"},
					token.Pos(1953879): {ID: "test-1", Loc: "test"},
					token.Pos(1953937): {ID: "test-2", Loc: "test"},
					token.Pos(1953992): {ID: "test-1", Loc: "test"},
					token.Pos(1954093): {ID: "test-1", Loc: "test"},
				},
			},
			stg: np,
			sts: map[string]gopium.Struct{
				"test-1": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"test-2": gopium.Struct{
					Name: "b",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/flat.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
			},
		},
		"flat struct pkg should visit only expected structs with regex": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`A`),
			m:   m,
			p:   data.NewParser("flat"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_flat_file.go:10": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:41": gopium.Struct{
					Name: "AZ",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/flat.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
		"nested structs pkg should visit only top level structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("nested"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_nested_file.go:7": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_nested_file.go:11": gopium.Struct{
					Name: "b",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/nested.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_nested_file.go:16": gopium.Struct{
					Name: "C",
					Fields: []gopium.Field{
						{
							Name:  "c",
							Type:  "[]string",
							Size:  24,
							Align: 8,
						},
						{
							Name:     "A",
							Type:     "struct{b 1pkg/gopium/tests/data/nested.b; z 1pkg/gopium/tests/data/nested.A}",
							Size:     24,
							Align:    8,
							Exported: true,
						},
					},
				},
				"tests_data_nested_file.go:63": gopium.Struct{
					Name: "Z",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:     "C",
							Type:     "1pkg/gopium/tests/data/nested.C",
							Size:     48,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
		"multi structs pkg should visit only expected top level structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`[AZ]`),
			m:   m,
			p:   data.NewParser("multi"),
			stg: pck,
			sts: map[string]gopium.Struct{
				"tests_data_multi_file-1.go:9": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_multi_file-3.go:17": gopium.Struct{
					Name: "AZ",
					Fields: []gopium.Field{
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/multi.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
				"tests_data_multi_file-3.go:27": gopium.Struct{
					Name: "Zeze",
					Fields: []gopium.Field{
						{
							Name:     "AZ",
							Type:     "1pkg/gopium/tests/data/multi.AZ",
							Size:     32,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/multi.D",
							Size:     24,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:     "AWA",
							Type:     "1pkg/gopium/tests/data/multi.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:     "ze",
							Type:     "1pkg/gopium/tests/data/multi.ze",
							Size:     16,
							Align:    8,
							Embedded: true,
						},
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			pkg, loc, err := tcase.p.ParseTypes(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to %v", err, nil)
			}
			ref := collections.NewReference(true)
			m := &maven{exp: m, loc: loc, ref: ref}
			if tcase.loc != nil {
				m.loc = tcase.loc
			}
			ch := make(appliedCh)
			// exec
			go vscope(tcase.ctx, pkg.Scope(), tcase.r, tcase.stg, m, ch)
			// check
			for applied := range ch {
				// if error occurred skip
				if applied.Err != nil {
					if !reflect.DeepEqual(applied.Err, tcase.err) {
						t.Errorf("actual %v doesn't equal to expected %v", applied.Err, tcase.err)
					}
					continue
				}
				// otherwise check all struct
				// against structs map
				if st, ok := tcase.sts[applied.ID]; ok {
					if !reflect.DeepEqual(applied.R, st) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", applied.ID, applied.R, st)
					}
					delete(tcase.sts, applied.ID)
				} else {
					t.Errorf("actual %v doesn't equal to expected %v", applied.ID, "")
				}
			}
			// check that map has been drained
			if !reflect.DeepEqual(tcase.sts, make(map[string]gopium.Struct)) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.sts, make(map[string]gopium.Struct))
			}
		})
	}
}

func TestVdeep(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	b := strategies.Builder{}
	np, err := b.Build(strategies.Nope)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pck, err := b.Build(strategies.Pack)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 64, 64, 64)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		ctx context.Context
		r   *regexp.Regexp
		m   gopium.Maven
		p   gopium.TypeParser
		loc gopium.Locator
		stg gopium.Strategy
		sts map[string]gopium.Struct
		err error
	}{
		"empty pkg should visit nothing": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("empty"),
			stg: np,
			sts: make(map[string]gopium.Struct),
		},
		"single struct pkg should visit the struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_single_file.go:5": gopium.Struct{
					Name: "Single",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "string",
							Size:     16,
							Align:    8,
							Exported: true,
						},
						{
							Name:     "B",
							Type:     "string",
							Size:     16,
							Align:    8,
							Exported: true,
						},
						{
							Name:     "C",
							Type:     "string",
							Size:     16,
							Align:    8,
							Exported: true,
						},
					},
				},
			},
		},
		"single struct pkg should visit nothing on canceled context": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
			stg: np,
			sts: make(map[string]gopium.Struct),
			err: context.Canceled,
		},
		"nested struct pkg should visit nothing on canceled context": {
			ctx: &mocks.Context{After: 2},
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("nested"),
			stg: np,
			sts: make(map[string]gopium.Struct),
			err: context.Canceled,
		},
		"flat struct pkg should visit all structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_flat_file.go:10": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:16": gopium.Struct{
					Name: "b",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/flat.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:21": gopium.Struct{
					Name: "C",
					Fields: []gopium.Field{
						{
							Name:  "c",
							Type:  "[]string",
							Size:  24,
							Align: 8,
						},
						{
							Name:     "A",
							Type:     "struct{b 1pkg/gopium/tests/data/flat.b; z 1pkg/gopium/tests/data/flat.A}",
							Size:     24,
							Align:    8,
							Exported: true,
						},
					},
				},
				"tests_data_flat_file.go:29": gopium.Struct{
					Name: "c1",
					Fields: []gopium.Field{
						{
							Name:  "c",
							Type:  "[]string",
							Size:  24,
							Align: 8,
						},
						{
							Name:     "A",
							Type:     "struct{b 1pkg/gopium/tests/data/flat.b; z 1pkg/gopium/tests/data/flat.A}",
							Size:     24,
							Align:    8,
							Exported: true,
						},
					},
				},
				"tests_data_flat_file.go:32": gopium.Struct{
					Name: "D",
					Fields: []gopium.Field{
						{
							Name:  "t",
							Type:  "[13]byte",
							Size:  13,
							Align: 1,
						},
						{
							Name:  "b",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:  "_",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:41": gopium.Struct{
					Name: "AZ",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/flat.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
		"flat struct pkg should visit the struct on same loc": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
			loc: mocks.Locator{
				Poses: map[token.Pos]mocks.Pos{
					token.Pos(1799681): {ID: "test-1", Loc: "test"},
					token.Pos(1799769): {ID: "test-2", Loc: "test"},
					token.Pos(1799802): {ID: "test-1", Loc: "test"},
					token.Pos(1799860): {ID: "test-2", Loc: "test"},
					token.Pos(1799915): {ID: "test-1", Loc: "test"},
					token.Pos(1800016): {ID: "test-1", Loc: "test"},

					token.Pos(1953758): {ID: "test-1", Loc: "test"},
					token.Pos(1953846): {ID: "test-2", Loc: "test"},
					token.Pos(1953879): {ID: "test-1", Loc: "test"},
					token.Pos(1953937): {ID: "test-2", Loc: "test"},
					token.Pos(1953992): {ID: "test-1", Loc: "test"},
					token.Pos(1954093): {ID: "test-1", Loc: "test"},
				},
			},
			stg: np,
			sts: map[string]gopium.Struct{
				"test-1": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"test-2": gopium.Struct{
					Name: "b",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/flat.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
			},
		},
		"flat struct pkg should visit only expected structs with regex": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`A`),
			m:   m,
			p:   data.NewParser("flat"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_flat_file.go:10": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_flat_file.go:41": gopium.Struct{
					Name: "AZ",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/flat.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
		"nested structs pkg should visit all levels structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("nested"),
			stg: np,
			sts: map[string]gopium.Struct{
				"tests_data_nested_file.go:7": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_nested_file.go:11": gopium.Struct{
					Name: "b",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/nested.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_nested_file.go:16": gopium.Struct{
					Name: "C",
					Fields: []gopium.Field{
						{
							Name:  "c",
							Type:  "[]string",
							Size:  24,
							Align: 8,
						},
						{
							Name:     "A",
							Type:     "struct{b 1pkg/gopium/tests/data/nested.b; z 1pkg/gopium/tests/data/nested.A}",
							Size:     24,
							Align:    8,
							Exported: true,
						},
					},
				},
				"tests_data_nested_file.go:25": gopium.Struct{
					Name: "B",
					Fields: []gopium.Field{
						{
							Name:     "b",
							Type:     "1pkg/gopium/tests/data/nested.b",
							Size:     16,
							Align:    8,
							Embedded: true,
						},
					},
				},
				"tests_data_nested_file.go:28": gopium.Struct{
					Name: "b1",
					Fields: []gopium.Field{
						{
							Name:     "A",
							Type:     "1pkg/gopium/tests/data/nested.A",
							Size:     8,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:  "b",
							Type:  "float64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_nested_file.go:37": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int32",
							Size:  4,
							Align: 4,
						},
					},
				},
				"tests_data_nested_file.go:40": gopium.Struct{
					Name: "a1",
					Fields: []gopium.Field{
						{
							Name:  "i",
							Type:  "interface{}",
							Size:  16,
							Align: 8,
						},
					},
				},
				"tests_data_nested_file.go:46": gopium.Struct{
					Name: "a1",
					Fields: []gopium.Field{
						{
							Name:  "i",
							Type:  "struct{}",
							Size:  0,
							Align: 1,
						},
					},
				},
				"tests_data_nested_file.go:63": gopium.Struct{
					Name: "Z",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:     "C",
							Type:     "1pkg/gopium/tests/data/nested.C",
							Size:     48,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
		"multi structs pkg should visit all expected levels structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`[AZ]`),
			m:   m,
			p:   data.NewParser("multi"),
			stg: pck,
			sts: map[string]gopium.Struct{
				"tests_data_multi_file-1.go:9": gopium.Struct{
					Name: "A",
					Fields: []gopium.Field{
						{
							Name:  "a",
							Type:  "int64",
							Size:  8,
							Align: 8,
						},
					},
				},
				"tests_data_multi_file-3.go:17": gopium.Struct{
					Name: "AZ",
					Fields: []gopium.Field{
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/multi.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
				"tests_data_multi_file-3.go:27": gopium.Struct{
					Name: "Zeze",
					Fields: []gopium.Field{
						{
							Name:     "AZ",
							Type:     "1pkg/gopium/tests/data/multi.AZ",
							Size:     32,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/multi.D",
							Size:     24,
							Align:    8,
							Exported: true,
							Embedded: true,
						},
						{
							Name:     "AWA",
							Type:     "1pkg/gopium/tests/data/multi.D",
							Size:     24,
							Align:    8,
							Exported: true,
						},
						{
							Name:     "ze",
							Type:     "1pkg/gopium/tests/data/multi.ze",
							Size:     16,
							Align:    8,
							Embedded: true,
						},
					},
				},
				"tests_data_multi_file-1.go:29": gopium.Struct{
					Name: "TestAZ",
					Fields: []gopium.Field{
						{
							Name:     "D",
							Type:     "1pkg/gopium/tests/data/multi.A",
							Size:     8,
							Align:    8,
							Exported: true,
						},
						{
							Name:  "a",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
						{
							Name:  "z",
							Type:  "bool",
							Size:  1,
							Align: 1,
						},
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			pkg, loc, err := tcase.p.ParseTypes(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to %v", err, nil)
			}
			ref := collections.NewReference(true)
			m := &maven{exp: m, loc: loc, ref: ref}
			if tcase.loc != nil {
				m.loc = tcase.loc
			}
			ch := make(appliedCh)
			// exec
			go vdeep(tcase.ctx, pkg.Scope(), tcase.r, tcase.stg, m, ch)
			// check
			for applied := range ch {
				// if error occurred skip
				if applied.Err != nil {
					if !reflect.DeepEqual(applied.Err, tcase.err) {
						t.Errorf("actual %v doesn't equal to expected %v", applied.Err, tcase.err)
					}
					continue
				}
				// otherwise check all struct
				// against structs map
				if st, ok := tcase.sts[applied.ID]; ok {
					if !reflect.DeepEqual(applied.R, st) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", applied.ID, applied.R, st)
					}
					delete(tcase.sts, applied.ID)
				} else {
					t.Errorf("actual %v doesn't equal to expected %v", applied.ID, "")
				}
			}
			// check that map has been drained
			if !reflect.DeepEqual(tcase.sts, make(map[string]gopium.Struct)) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.sts, make(map[string]gopium.Struct))
			}
		})
	}
}
