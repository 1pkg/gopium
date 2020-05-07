package walkers

import (
	"context"
	"errors"
	"go/build"
	"go/types"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/fmtio"
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
	p := fmtio.Goprint(0, 4, false)
	table := map[string]struct {
		ctx  context.Context
		r    *regexp.Regexp
		p    gopium.Parser
		a    astutil.Apply
		w    fmtio.Writer
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
			a:   astutil.FFN,
			stg: np,
			sts: map[string][]byte{
				"/src/1pkg/gopium/tests/data/empty/file.go": []byte(`
//+build tests_data

package empty
`),
			},
		},
		"single struct pkg should visit the struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.FFN,
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
		"single struct pkg should visit nothing on canceled context": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.FFN,
			stg: np,
			sts: make(map[string][]byte),
			err: cctx.Err(),
		},
		"single struct pkg should visit nothing on type parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   mocks.Parser{Terr: errors.New("test-1")},
			a:   astutil.FFN,
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
			a:   astutil.FFN,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-2"),
		},
		"single struct pkg should visit nothing on strategy error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.FFN,
			stg: mocks.Strategy{Err: errors.New("test-3")},
			sts: make(map[string][]byte),
			err: errors.New("test-3"),
		},
		"single struct pkg should visit nothing on persist error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.FFN,
			w:   (&mocks.Writer{Err: errors.New("test-4")}).Writer,
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
		"multi structs pkg should visit all expected levels structs with deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`(A|Z)`),
			p:    data.NewParser("multi"),
			a:    astutil.FFN,
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
		"multi structs pkg should visit all expected levels structs without deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`(A|Z)`),
			p:    data.NewParser("multi"),
			a:    astutil.FFN,
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
			// prepare
			w := &mocks.Writer{}
			wast := wast{
				apply:  tcase.a,
				writer: w.Writer,
			}.With(tcase.p, m, p, tcase.deep, tcase.bref)
			if tcase.w != nil {
				wast.writer = tcase.w
			}
			// exec
			err := wast.Visit(tcase.ctx, tcase.r, tcase.stg)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			for id, buf := range w.Buffers {
				// remove gopath from id
				id = strings.Replace(id, build.Default.GOPATH, "", 1)
				// check all struct
				// against bytes map
				if st, ok := tcase.sts[id]; ok {
					// format actual and expected identically
					actual := strings.Trim(string(buf.Bytes()), "\n")
					expected := strings.Trim(string(st), "\n")
					if !reflect.DeepEqual(actual, expected) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", id, actual, expected)
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
