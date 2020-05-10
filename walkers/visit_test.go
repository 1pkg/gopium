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
			stg:  mocks.Strategy{},
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
			stg:  mocks.Strategy{},
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
			stg:  mocks.Strategy{},
			ch:   make(appliedCh),
			deep: false,
			ctx:  context.Background(),
			s:    &types.Scope{},
		},
		"with visit should return expected govisit func without all flags": {
			exp: mocks.Maven{},
			loc: mocks.Locator{},
			r:   regexp.MustCompile(`.*`),
			stg: mocks.Strategy{},
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
			if gvisit == nil {
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
				"5-b0652be9c761c2f34deff8a560333dd372ee062bb1dbcba6a79647fdc3205919": gopium.Struct{
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
		"flat struct pkg should visit all structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
			stg: np,
			sts: map[string]gopium.Struct{
				"10-cc533a18fa665ce942eb8127f87a8e3f1f007bc921cd29d5c731442351f9cb1f": gopium.Struct{
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
				"16-00383e272efed5ebb4ed09e9a5a5d1ac6c5c66ab722d1b5aabdbe6be239b1b68": gopium.Struct{
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
				"21-5f1ccb6e883ad93653d67eeaa568de2691fd098f873ba4b89699b1127eb9368f": gopium.Struct{
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
				"29-f8b876109e914c49453aa46663a932d2e0227265423b9053f92e45b0df397228": gopium.Struct{
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
				"32-6de12ac7df310c06ecb758c3b0f101240494266e7449a04a1118fcffb1f5e7ed": gopium.Struct{
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
				"41-849b383f3d6ca6222a423e60766e81d10a536633dd407fed11fab9cead6f43d5": gopium.Struct{
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
				"10-cc533a18fa665ce942eb8127f87a8e3f1f007bc921cd29d5c731442351f9cb1f": gopium.Struct{
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
				"41-849b383f3d6ca6222a423e60766e81d10a536633dd407fed11fab9cead6f43d5": gopium.Struct{
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
				"7-d68df23f11155cc8dc251a831180fa4ea7a0632b9ad7da370c767ee439ca965a": gopium.Struct{
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
				"11-873eeacc02d5a57f16b51905dfffd8b4d6696e413f47e31b3ad3a7d6c6d9a80a": gopium.Struct{
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
				"16-5f1da0f40f8cc9353e61f010511ebfe91a33c578b115af07d05596c9180db986": gopium.Struct{
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
				"63-9e99e5c6b447375ce9a92d470e79e76dbc4e690c8f71cc57ca19a14b14b43e19": gopium.Struct{
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
			r:   regexp.MustCompile(`(A|Z)`),
			m:   m,
			p:   data.NewParser("multi"),
			stg: pck,
			sts: map[string]gopium.Struct{
				"9-7d858286ee3f6bdbb9c740b5333435af40ec918bdeec00ececacf5ab9764f09b": gopium.Struct{
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
				"17-342e1133d9f044ad74cd048f681aad0efcca3407b8fe3b972c96eb92d034fd04": gopium.Struct{
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
				"27-6a3c1ba2a278b9b24c0d76ad232bba0f0b0abd806f9cbb6e0910966f761e5130": gopium.Struct{
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
			ref := collections.NewReference(false)
			m := &maven{exp: m, loc: loc, ref: ref}
			if tcase.loc != nil {
				m.loc = tcase.loc
			}
			ch := make(appliedCh)
			// exec
			go vscope(tcase.ctx, pkg.Scope(), tcase.r, tcase.stg, m, ch)
			// check
			for applied := range ch {
				// if error occured skip
				if applied.Error != nil {
					if !reflect.DeepEqual(applied.Error, tcase.err) {
						t.Errorf("actual %v doesn't equal to expected %v", applied.Error, tcase.err)
					}
					continue
				}
				// otherwise check all struct
				// against structs map
				if st, ok := tcase.sts[applied.ID]; ok {
					if !reflect.DeepEqual(applied.Result, st) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", applied.ID, applied.Result, st)
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
				"5-b0652be9c761c2f34deff8a560333dd372ee062bb1dbcba6a79647fdc3205919": gopium.Struct{
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
		"flat struct pkg should visit all structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
			stg: np,
			sts: map[string]gopium.Struct{
				"10-cc533a18fa665ce942eb8127f87a8e3f1f007bc921cd29d5c731442351f9cb1f": gopium.Struct{
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
				"16-00383e272efed5ebb4ed09e9a5a5d1ac6c5c66ab722d1b5aabdbe6be239b1b68": gopium.Struct{
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
				"21-5f1ccb6e883ad93653d67eeaa568de2691fd098f873ba4b89699b1127eb9368f": gopium.Struct{
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
				"29-f8b876109e914c49453aa46663a932d2e0227265423b9053f92e45b0df397228": gopium.Struct{
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
				"32-6de12ac7df310c06ecb758c3b0f101240494266e7449a04a1118fcffb1f5e7ed": gopium.Struct{
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
				"41-849b383f3d6ca6222a423e60766e81d10a536633dd407fed11fab9cead6f43d5": gopium.Struct{
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
				"10-cc533a18fa665ce942eb8127f87a8e3f1f007bc921cd29d5c731442351f9cb1f": gopium.Struct{
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
				"41-849b383f3d6ca6222a423e60766e81d10a536633dd407fed11fab9cead6f43d5": gopium.Struct{
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
				"7-d68df23f11155cc8dc251a831180fa4ea7a0632b9ad7da370c767ee439ca965a": gopium.Struct{
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
				"11-873eeacc02d5a57f16b51905dfffd8b4d6696e413f47e31b3ad3a7d6c6d9a80a": gopium.Struct{
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
				"16-5f1da0f40f8cc9353e61f010511ebfe91a33c578b115af07d05596c9180db986": gopium.Struct{
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
				"25-ee3b039ceabb27a11b329533e48f1db7207ea4b2bfd174bd286524237da7bc7c": gopium.Struct{
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
				"28-a058328567dbcd6c2658f94dde6c835d5560f7a25e2043e260c3423116f931a7": gopium.Struct{
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
				"37-9deae08ffbbacf4700f78570260603d4200a04d8bc7011f442fde7a82b420bb9": gopium.Struct{
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
				"40-7bab018214c27a1b6d1c5cdf37adeb2cd468b28723dd13e9eb3811cb41e1b725": gopium.Struct{
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
				"46-3054b137d9628f8dda52acdb8084e593cf42cc3cbcf89aa454eb3cd2c240e593": gopium.Struct{
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
				"63-9e99e5c6b447375ce9a92d470e79e76dbc4e690c8f71cc57ca19a14b14b43e19": gopium.Struct{
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
			r:   regexp.MustCompile(`(A|Z)`),
			m:   m,
			p:   data.NewParser("multi"),
			stg: pck,
			sts: map[string]gopium.Struct{
				"9-7d858286ee3f6bdbb9c740b5333435af40ec918bdeec00ececacf5ab9764f09b": gopium.Struct{
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
				"17-342e1133d9f044ad74cd048f681aad0efcca3407b8fe3b972c96eb92d034fd04": gopium.Struct{
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
				"27-6a3c1ba2a278b9b24c0d76ad232bba0f0b0abd806f9cbb6e0910966f761e5130": gopium.Struct{
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
				"29-6dc854454cff4b7c6b7ba90ba55fa564c21409c5a107cf402dd2e582d44dd32a": gopium.Struct{
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
			ref := collections.NewReference(false)
			m := &maven{exp: m, loc: loc, ref: ref}
			if tcase.loc != nil {
				m.loc = tcase.loc
			}
			ch := make(appliedCh)
			// exec
			go vdeep(tcase.ctx, pkg.Scope(), tcase.r, tcase.stg, m, ch)
			// check
			for applied := range ch {
				// if error occured skip
				if applied.Error != nil {
					if !reflect.DeepEqual(applied.Error, tcase.err) {
						t.Errorf("actual %v doesn't equal to expected %v", applied.Error, tcase.err)
					}
					continue
				}
				// otherwise check all struct
				// against structs map
				if st, ok := tcase.sts[applied.ID]; ok {
					if !reflect.DeepEqual(applied.Result, st) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", applied.ID, applied.Result, st)
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
