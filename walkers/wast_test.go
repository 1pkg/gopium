package walkers

import (
	"context"
	"errors"
	"go/types"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/astutil/apply"
	"1pkg/gopium/astutil/print"
	"1pkg/gopium/strategies"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
	"1pkg/gopium/typepkg"
)

func TestWast(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	b := strategies.Builder{}
	np, err := b.Build(strategies.Nope)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pck, err := b.Build(strategies.Pack)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 64, 64, 64)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	p := print.GoPrinter(0, 4, false)
	table := map[string]struct {
		ctx  context.Context
		r    *regexp.Regexp
		p    gopium.Parser
		a    astutil.Apply
		prs  astutil.Persist
		stg  gopium.Strategy
		deep bool
		bref bool
		sts  map[string][]byte
		err  error
	}{
		"empty pkg should visit nothing": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("empty"),
			a:   apply.SFN,
			stg: np,
			sts: map[string][]byte{
				"/src/1pkg/gopium/tests/data/empty/file.go": []byte(`
//+build tests_data

package empty
`),
			},
		},
		"single struct pkg should visit the single struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   apply.SFN,
			stg: np,
			sts: map[string][]byte{
				"/src/1pkg/gopium/tests/data/single/file.go": []byte(`
//+build tests_data

package single

type Single struct {
	A	string
	B	string
	C	string
}
`),
			},
		},
		"single struct pkg should visit nothing on context cancelation": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   apply.SFN,
			stg: np,
			sts: make(map[string][]byte),
			err: cctx.Err(),
		},
		"single struct pkg should visit nothing on type parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   mocks.Parser{Terr: errors.New("test-1")},
			a:   apply.SFN,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-1"),
		},
		"single struct pkg should visit nothing on ast parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p: mocks.Parser{
				Types: types.NewPackage("", ""),
				Aerr:  errors.New("test-2"),
			},
			a:   apply.SFN,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-2"),
		},
		"single struct pkg should visit nothing on strategy error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   apply.SFN,
			stg: mocks.Strategy{Err: errors.New("test-3")},
			sts: make(map[string][]byte),
			err: errors.New("test-3"),
		},
		"single struct pkg should visit nothing on persist error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   apply.SFN,
			prs: (&mocks.Persist{Err: errors.New("test-4")}).Persist,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-4"),
		},
		"single struct pkg should visit nothing on apply error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   (&mocks.Apply{Err: errors.New("test-5")}).Apply,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-5"),
		},
		"multi structs pkg should visit all relevant levels structs with deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`(A|Z)`),
			p:    data.NewParser("multi"),
			a:    apply.SFN,
			stg:  pck,
			deep: true,
			sts: map[string][]byte{
				"/src/1pkg/gopium/tests/data/multi/file-1.go": []byte(`
//+build tests_data

package multi

import (
	"strings"
)

type A struct {
	a int64
}

var a1 string = strings.Join([]string{"a", "b", "c"}, "|")

type b struct {
	A
	b	float64
}

type C struct {
	c	[]string
	A	struct {
		b	b
		z	A
	}
}

func scope() {
	type TestAZ struct {
		D	A
		a	bool
		z	bool
	}
}
`),
				"/src/1pkg/gopium/tests/data/multi/file-2.go": []byte(`
//+build tests_data

package multi

import "errors"

func scope1() error {
	type B struct {
		b
	}
	type b1 b
	type b2 struct {
		A
		b	float64
	}
	return errors.New("test data")
}
`),
				"/src/1pkg/gopium/tests/data/multi/file-3.go": []byte(`
//+build tests_data

package multi

type c1 C

// table := []struct{A string}{{A: "test"}}
type D struct {
	t	[13]byte
	b	bool
	_	int64
}

/* ggg := func (interface{}){} */
type AW func() error

type AZ struct {
	D	D
	a	bool
	z	bool
}

type ze interface {
	AW() AW
}

type Zeze struct {
	ze
	D
	AZ
	AWA	D
}

// test comment
type (
	d1	int64
	d2	float64
	d3	string
)
`),
			},
		},
		"multi structs pkg should visit all relevant levels structs without deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`(A|Z)`),
			p:    data.NewParser("multi"),
			a:    apply.SFN,
			stg:  pck,
			bref: true,
			sts: map[string][]byte{
				"/src/1pkg/gopium/tests/data/multi/file-1.go": []byte(`
//+build tests_data

package multi

import (
	"strings"
)

type A struct {
	a int64
}

var a1 string = strings.Join([]string{"a", "b", "c"}, "|")

type b struct {
	A
	b	float64
}

type C struct {
	c	[]string
	A	struct {
		b	b
		z	A
	}
}

func scope() {
	type TestAZ struct {
		a	bool
		D	A
		z	bool
	}
}
`),
				"/src/1pkg/gopium/tests/data/multi/file-2.go": []byte(`
//+build tests_data

package multi

import "errors"

func scope1() error {
	type B struct {
		b
	}
	type b1 b
	type b2 struct {
		A
		b	float64
	}
	return errors.New("test data")
}
`),
				"/src/1pkg/gopium/tests/data/multi/file-3.go": []byte(`
//+build tests_data

package multi

type c1 C

// table := []struct{A string}{{A: "test"}}
type D struct {
	t	[13]byte
	b	bool
	_	int64
}

/* ggg := func (interface{}){} */
type AW func() error

type AZ struct {
	D	D
	a	bool
	z	bool
}

type ze interface {
	AW() AW
}

type Zeze struct {
	ze
	D
	AZ
	AWA	D
}

// test comment
type (
	d1	int64
	d2	float64
	d3	string
)
`),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			prs := &mocks.Persist{}
			wast := wast{
				apply:   tcase.a,
				persist: prs.Persist,
			}.With(tcase.p, m, p, tcase.deep, tcase.bref)
			if tcase.prs != nil {
				wast.persist = tcase.prs
			}
			err := wast.Visit(tcase.ctx, tcase.r, tcase.stg)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			for id, buf := range prs.Buffers {
				// check all struct
				// against bytes map
				if st, ok := tcase.sts[id]; ok {
					// format actual and expected identically
					stract, strexp := strings.Trim(string(buf.Bytes()), "\n"), strings.Trim(string(st), "\n")
					if !reflect.DeepEqual(stract, strexp) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", id, stract, strexp)
					}
					delete(tcase.sts, id)
				} else {
					t.Errorf("actual %v doesn't equal to expected %v", id, "")
				}
			}
			// check that map has been drained
			if !reflect.DeepEqual(tcase.sts, make(map[string][]byte)) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.sts, make(map[string][]byte))
			}
		})
	}
}
