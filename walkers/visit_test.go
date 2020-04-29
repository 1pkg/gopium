package walkers

import (
	"context"
	"go/types"
	"reflect"
	"regexp"
	"sync"
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
		exp   gopium.Exposer
		loc   gopium.Locator
		bref  bool
		regex *regexp.Regexp
		stg   gopium.Strategy
		ch    appliedCh
		deep  bool
		ctx   context.Context
		s     *types.Scope
	}{
		"with visit should create valid govisit func": {
			exp:   mocks.Maven{},
			loc:   mocks.Locator{},
			bref:  true,
			regex: regexp.MustCompile(`.*`),
			stg:   mocks.Strategy{},
			ch:    make(appliedCh),
			deep:  true,
			ctx:   context.Background(),
			s:     &types.Scope{},
		},
		"with visit should create valid govisit func even without bref flag": {
			exp:   mocks.Maven{},
			loc:   mocks.Locator{},
			bref:  false,
			regex: regexp.MustCompile(`.*`),
			stg:   mocks.Strategy{},
			ch:    make(appliedCh),
			deep:  true,
			ctx:   context.Background(),
			s:     &types.Scope{},
		},
		"with visit should create valid govisit func even without deep flag": {
			exp:   mocks.Maven{},
			loc:   mocks.Locator{},
			bref:  true,
			regex: regexp.MustCompile(`.*`),
			stg:   mocks.Strategy{},
			ch:    make(appliedCh),
			deep:  false,
			ctx:   context.Background(),
			s:     &types.Scope{},
		},
		"with visit should create valid govisit func even without flags": {
			exp:   mocks.Maven{},
			loc:   mocks.Locator{},
			regex: regexp.MustCompile(`.*`),
			stg:   mocks.Strategy{},
			ch:    make(appliedCh),
			ctx:   context.Background(),
			s:     &types.Scope{},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			gvisit := with(tcase.exp, tcase.loc, tcase.bref).
				visit(tcase.regex, tcase.stg, tcase.ch, tcase.deep)
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
	var wg sync.WaitGroup
	stg, err := strategies.Builder{}.Build(strategies.Nope)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 64, 64, 64)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		ctx context.Context
		r   *regexp.Regexp
		m   gopium.Maven
		p   gopium.TypeParser
		sts map[string]gopium.Struct
		err error
	}{
		"empty pkg should visit nothing": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("empty"),
			sts: make(map[string]gopium.Struct),
		},
		"single struct pkg should visit the single struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
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
		"single struct pkg should visit nothing on context cancelation": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("single"),
			sts: make(map[string]gopium.Struct),
			err: cctx.Err(),
		},
		"flat struct pkg should visit all structs": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			m:   m,
			p:   data.NewParser("flat"),
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
		"flat struct pkg should visit only relevant structs with regex": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`A`),
			m:   m,
			p:   data.NewParser("flat"),
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
				"68-5c6eddb8c98e95eca87eead69e3d84ebf9215ad81246ec2622d05ed19b2e0f71": gopium.Struct{
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
	}
	for name, tcase := range table {
		// run all parser tests
		// in separate goroutine
		wg.Add(1)
		name := name
		tcase := tcase
		go func(t *testing.T) {
			defer wg.Done()
			t.Run(name, func(t *testing.T) {
				// exec
				pkg, loc, err := tcase.p.ParseTypes(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				ref := collections.NewReference(false)
				m := &maven{exp: m, loc: loc, ref: ref}
				ch := make(appliedCh)
				// check
				go vscope(tcase.ctx, pkg.Scope(), tcase.r, stg, m, ch)
				for applied := range ch {
					// if error occured check it
					if applied.Error != nil {
						if !reflect.DeepEqual(applied.Error, tcase.err) {
							t.Errorf("actual %v doesn't equal to expected %v", applied.Error, tcase.err)
						}
						return
					}
					// otherwise check all struct
					// against structs map
					if st, ok := tcase.sts[applied.ID]; ok {
						if !reflect.DeepEqual(applied.Result, st) {
							t.Errorf("actual %v doesn't equal to expected %v", applied.Result, st)
						}
						delete(tcase.sts, applied.ID)
					} else {
						t.Errorf("actual %v doesn't equal to expected %v", applied.ID, "")
					}
				}
				// check that map has been drained
				if stsl := len(tcase.sts); stsl > 0 {
					t.Errorf("actual %v doesn't equal to expected %v", stsl, 0)
				}
			})
		}(t)
	}
	// wait util tests finish
	wg.Wait()
}
