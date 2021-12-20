package walkers

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/1pkg/gopium/fmtio"
	"github.com/1pkg/gopium/fmtio/astutil"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/strategies"
	"github.com/1pkg/gopium/tests/data"
	"github.com/1pkg/gopium/tests/mocks"
	"github.com/1pkg/gopium/typepkg"
)

func TestWast(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	b := strategies.Builder{}
	np, err := b.Build(strategies.Ignore)
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
	p := fmtio.Gofmt{}
	table := map[string]struct {
		ctx  context.Context
		r    *regexp.Regexp
		p    gopium.Parser
		a    gopium.Apply
		sp   gopium.Persister
		w    gopium.CategoryWriter
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
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{},
		},
		"single struct pkg should visit the struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{
				"tests_data_single_file.go": []byte(`
//go:build tests_data

package single

type Single struct {
	A string
	B string
	C string
}
`),
			},
		},
		"single struct pkg should visit nothing on canceled context": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{},
			err: context.Canceled,
		},
		"single struct pkg should visit nothing on type parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   mocks.Parser{Typeserr: errors.New("test-1")},
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-1"),
		},
		"single struct pkg should visit nothing on ast parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   mocks.Parser{Parser: data.NewParser("single"), Asterr: errors.New("test-2")},
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-2"),
		},
		"single struct pkg should visit nothing on strategy error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: &mocks.Strategy{Err: errors.New("test-3")},
			sts: map[string][]byte{},
			err: errors.New("test-3"),
		},
		"single struct pkg should visit nothing on persist error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: (&mocks.Writer{Gerr: errors.New("test-4")})},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-4"),
		},
		"single struct pkg should visit nothing on cat persist error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.UFFN,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: (&mocks.Writer{Cerr: errors.New("test-5")})},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-5"),
		},
		"single struct pkg should visit nothing on apply error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   (&mocks.Apply{Err: errors.New("test-6")}).Apply,
			sp:  astutil.Package{},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-6"),
		},
		"single struct pkg should visit nothing on persister error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			a:   astutil.UFFN,
			sp:  mocks.Persister{Err: errors.New("test-7")},
			w:   data.Writer{Writer: &mocks.Writer{}},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-7"),
		},
		"multi structs pkg should visit all expected levels structs with deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`([AZ])`),
			p:    data.NewParser("multi"),
			a:    astutil.UFFN,
			sp:   astutil.Package{},
			w:    data.Writer{Writer: &mocks.Writer{}},
			stg:  pck,
			deep: true,
			sts: map[string][]byte{
				"tests_data_multi_file-1.go": []byte(`
//go:build tests_data

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
	b float64
}

type C struct {
	c []string
	A struct {
		b b
		z A
	}
}

func scope() {
	type TestAZ struct {
		D A
		a bool
		z bool
	}
}
`),
				"tests_data_multi_file-3.go": []byte(`
//go:build tests_data

package multi

type c1 C

// table := []struct{A string}{{A: "test"}}
type D struct {
	t [13]byte
	b bool
	_ int64
}

/* ggg := func (interface{}){} */
type AW func() error

type AZ struct {
	D D
	a bool
	z bool
}

type ze interface {
	AW() AW
}

type Zeze struct {
	ze

	AZ
	D
	AWA D
}

// test comment
type (
	d1 int64
	d2 float64
	d3 string
)
`),
			},
		},
		"multi structs pkg should visit all expected levels structs without deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`([AZ])`),
			p:    data.NewParser("multi"),
			a:    astutil.UFFN,
			sp:   astutil.Package{},
			w:    data.Writer{Writer: &mocks.Writer{}},
			stg:  pck,
			bref: true,
			sts: map[string][]byte{
				"tests_data_multi_file-1.go": []byte(`
//go:build tests_data

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
	b float64
}

type C struct {
	c []string
	A struct {
		b b
		z A
	}
}

func scope() {
	type TestAZ struct {
		a bool
		D A
		z bool
	}
}
`),
				"tests_data_multi_file-3.go": []byte(`
//go:build tests_data

package multi

type c1 C

// table := []struct{A string}{{A: "test"}}
type D struct {
	t [13]byte
	b bool
	_ int64
}

/* ggg := func (interface{}){} */
type AW func() error

type AZ struct {
	D D
	a bool
	z bool
}

type ze interface {
	AW() AW
}

type Zeze struct {
	ze

	AZ
	D
	AWA D
}

// test comment
type (
	d1 int64
	d2 float64
	d3 string
)
`),
			},
		},
		"embedded structs pkg should visit all expected levels structs without deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`.*`),
			p:    data.NewParser("embedded"),
			a:    astutil.UFFN,
			sp:   astutil.Package{},
			w:    data.Writer{Writer: &mocks.Writer{}},
			stg:  pck,
			bref: true,
			sts: map[string][]byte{
				"tests_data_embedded_file.go": []byte(`
//go:build tests_data

package embedded

import "time"

type MetaLabaratory struct {
}

type Person struct {
	Birtday time.Time 'json:"birthday" db:"birthday"'
	Weight  float64   'json:"weight" db:"weight"'
	Height  float64   'json:"height" db:"height"'
}

type PatientObject struct {
	MetaLabaratory
	Person
	ID           string  'json:"id" db:"id"'
	Gender       string  'json:"gender" db:"gender"'
	PhoneNumber  *string 'json:"phone_number" db:"phone_number"'
	Email        *string 'json:"email" db:"email"'
	AddressTitle *string 'json:"address_title" db:"address_title"'
	Enrolled     bool    'json:"enrolled" db:"enrolled"'
}
`),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			wast := wast{
				apply:     tcase.a,
				persister: tcase.sp,
				writer:    tcase.w,
			}.With(tcase.p, m, p, tcase.deep, tcase.bref)
			// exec
			err := wast.Visit(tcase.ctx, tcase.r, tcase.stg)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// process checks only on success
			if tcase.err == nil {
				w := (tcase.w.(data.Writer)).Writer.(*mocks.Writer)
				for id, rwc := range w.RWCs {
					// check all struct
					// against bytes map
					if st, ok := tcase.sts[id]; ok {
						// read rwc to buffer
						var buf bytes.Buffer
						_, err := buf.ReadFrom(rwc)
						if !reflect.DeepEqual(err, nil) {
							t.Errorf("actual %v doesn't equal to expected %v", err, nil)
						}
						// format actual and expected identically
						actual := strings.Trim(buf.String(), "\n")
						expected := strings.ReplaceAll(strings.Trim(string(st), "\n"), "'", "`")
						if !reflect.DeepEqual(actual, expected) {
							t.Errorf("id %v actual %v doesn't equal to expected %v", id, actual, expected)
						}
						delete(tcase.sts, id)
					} else {
						t.Errorf("actual %v doesn't equal to expected %v", id, "")
					}
				}
				// check that map has been drained
				if !reflect.DeepEqual(tcase.sts, map[string][]byte{}) {
					t.Errorf("actual %v doesn't equal to expected %v", tcase.sts, map[string][]byte{})
				}
			}
		})
	}
}
