package fmtio

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
)

func TestPrinter(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p    gopium.Parser
		pr   Printer
		w    gopium.Writer
		skip bool
		ctx  context.Context
		r    map[string][]byte
		err  error
	}{
		"empty pkg should print nothing": {
			p:   data.NewParser("empty"),
			pr:  Goprint(0, 4, true),
			ctx: context.Background(),
			r: map[string][]byte{
				"tests_data_empty_file.go": []byte(`
//+build tests_data

package empty
`),
			},
		},
		"single struct pkg should print the struct": {
			p:   data.NewParser("single"),
			pr:  Goprint(0, 4, false),
			ctx: context.Background(),
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
		"single struct pkg should print nothing on canceled context": {
			p:   data.NewParser("single"),
			pr:  Goprint(0, 4, false),
			ctx: cctx,
			r:   map[string][]byte{},
			err: context.Canceled,
		},
		"single struct pkg should print nothing on canceled context in printer": {
			p:    data.NewParser("single"),
			pr:   Goprint(0, 4, false),
			skip: true,
			ctx:  cctx,
			r:    map[string][]byte{},
			err:  context.Canceled,
		},
		"single struct pkg should print nothing on persist error": {
			p:   data.NewParser("single"),
			pr:  Goprint(0, 4, false),
			ctx: context.Background(),
			w:   (&mocks.Writer{Gerr: errors.New("test-1")}),
			r:   map[string][]byte{},
			err: errors.New("test-1"),
		},
		"single struct pkg should print nothing on printer error": {
			p:   data.NewParser("single"),
			pr:  mocks.Printer{Err: errors.New("test-2")}.Printer,
			ctx: context.Background(),
			r: map[string][]byte{
				"tests_data_single_file.go": []byte(``),
			},
			err: errors.New("test-2"),
		},
		"multi structs pkg should print all expected levels structs": {
			p:   data.NewParser("multi"),
			pr:  Goprint(0, 4, false),
			ctx: context.Background(),
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
			writer := tcase.w
			w := &mocks.Writer{}
			if writer == nil {
				writer = w
			}
			pkg, loc, err := tcase.p.ParseAst(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
			}
			// exec
			if tcase.skip {
				err = tcase.pr(tcase.ctx, &bytes.Buffer{}, loc.Root(), pkg)
			} else {
				err = tcase.pr.Save(writer)(tcase.ctx, pkg, loc)
			}
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
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
		})
	}
}
