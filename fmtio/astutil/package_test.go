package astutil

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
)

func TestPackage(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		xp  gopium.Parser
		ctx context.Context
		p   gopium.Printer
		w   gopium.Writer
		r   map[string][]byte
		err error
	}{
		"empty pkg should print nothing": {
			xp:  data.NewParser("empty"),
			ctx: context.Background(),
			p:   fmtio.NewGoprinter(0, 4, true),
			w:   &mocks.Writer{},
			r: map[string][]byte{
				"tests_data_empty_file.go": []byte(`
//+build tests_data

package empty
`),
			},
		},
		"single struct pkg should print and persists the struct": {
			xp:  data.NewParser("single"),
			ctx: context.Background(),
			p:   fmtio.NewGoprinter(0, 4, false),
			w:   &mocks.Writer{},
			r: map[string][]byte{
				"tests_data_single_file.go": []byte(`
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
		"single struct pkg should persist nothing on canceled context": {
			xp:  data.NewParser("single"),
			ctx: cctx,
			p:   fmtio.NewGoprinter(0, 4, false),
			w:   &mocks.Writer{},
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"single struct pkg should persist nothing on persist error": {
			xp:  data.NewParser("single"),
			ctx: context.Background(),
			p:   fmtio.NewGoprinter(0, 4, false),
			w:   (&mocks.Writer{Gerr: errors.New("test-1")}),
			r:   map[string][]byte{},
			err: errors.New("test-1"),
		},
		"single struct pkg should persist nothing on printer error": {
			xp:  data.NewParser("single"),
			ctx: context.Background(),
			p:   mocks.Printer{Err: errors.New("test-2")},
			w:   &mocks.Writer{},
			r: map[string][]byte{
				"tests_data_single_file.go": []byte(``),
			},
			err: errors.New("test-2"),
		},
		"multi structs pkg should persist all expected levels structs": {
			xp:  data.NewParser("multi"),
			ctx: context.Background(),
			p:   fmtio.NewGoprinter(0, 4, false),
			w:   &mocks.Writer{},
			r: map[string][]byte{
				"tests_data_multi_file-1.go": []byte(`
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
				"tests_data_multi_file-2.go": []byte(`
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
				"tests_data_multi_file-3.go": []byte(`
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
	a	bool
	D	D
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
			pkg, loc, err := tcase.xp.ParseAst(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
			}
			// exec
			err = Package{}.Persist(tcase.ctx, tcase.p, tcase.w, loc, pkg)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// process checks only on success
			if tcase.err == nil {
				w := tcase.w.(*mocks.Writer)
				for name, rwc := range w.RWCs {
					// check all struct
					// against bytes map
					if st, ok := tcase.r[name]; ok {
						// read rwc to buffer
						var buf bytes.Buffer
						_, err := buf.ReadFrom(rwc)
						if !reflect.DeepEqual(err, nil) {
							t.Errorf("actual %v doesn't equal to expected %v", err, nil)
						}
						// format actual and expected identically
						actual := strings.Trim(string(buf.Bytes()), "\n")
						expected := strings.Trim(string(st), "\n")
						if !reflect.DeepEqual(actual, expected) {
							t.Errorf("name %v actual %v doesn't equal to expected %v", name, actual, expected)
						}
						delete(tcase.r, name)
					} else {
						t.Errorf("actual %v doesn't equal to expected %v", name, "")
					}
				}
				// check that map has been drained
				if !reflect.DeepEqual(tcase.r, map[string][]byte{}) {
					t.Errorf("actual %v doesn't equal to expected %v", tcase.r, map[string][]byte{})
				}
			}
		})
	}
}
