package fmtio

import (
	"context"
	"errors"
	"go/ast"
	"go/build"
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
		pr  Printer
		w   Writer
		ctx context.Context
		pkg *ast.Package
		p   gopium.Parser
		r   map[string][]byte
		err error
	}{
		"empty pkg should visit nothing": {
			pr:  Goprint(0, 4, true),
			ctx: context.Background(),
			p:   data.NewParser("empty"),
			r: map[string][]byte{
				"/src/1pkg/gopium/tests/data/empty/file.go": []byte(`
//+build tests_data

package empty
`),
			},
		},
		"single struct pkg should visit the single struct": {
			pr:  Goprint(0, 4, false),
			ctx: context.Background(),
			p:   data.NewParser("single"),
			r: map[string][]byte{
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
			pr:  Goprint(0, 4, false),
			ctx: cctx,
			p:   data.NewParser("single"),
			r:   make(map[string][]byte),
			err: cctx.Err(),
		},
		"single struct pkg should visit nothing on persist error": {
			pr:  Goprint(0, 4, false),
			ctx: context.Background(),
			p:   data.NewParser("single"),
			w:   (&mocks.Writer{Err: errors.New("test-1")}).Writer,
			r:   make(map[string][]byte),
			err: errors.New("test-1"),
		},
		"single struct pkg should visit nothing on printer error": {
			pr:  mocks.Printer{Err: errors.New("test-2")}.Printer,
			ctx: context.Background(),
			p:   data.NewParser("single"),
			r: map[string][]byte{
				"/src/1pkg/gopium/tests/data/single/file.go": []byte(``),
			},
			err: errors.New("test-2"),
		},
		"multi structs pkg should visit all relevant levels structs with deep": {
			pr:  Goprint(0, 4, false),
			ctx: context.Background(),
			p:   data.NewParser("multi"),
			r: map[string][]byte{
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
			// exec
			writer := tcase.w
			w := &mocks.Writer{}
			if writer == nil {
				writer = w.Writer
			}
			pkg, loc, err := tcase.p.ParseAst(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
			}
			err = tcase.pr.Save(writer)(tcase.ctx, pkg, loc)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			for name, buf := range w.Buffers {
				// remove gopath from collected id
				name = strings.Replace(name, build.Default.GOPATH, "", 1)
				// check all struct
				// against bytes map
				if st, ok := tcase.r[name]; ok {
					// format actual and expected identically
					stract, strexp := strings.Trim(string(buf.Bytes()), "\n"), strings.Trim(string(st), "\n")
					if !reflect.DeepEqual(stract, strexp) {
						t.Errorf("name %v actual %v doesn't equal to expected %v", name, stract, strexp)
					}
					delete(tcase.r, name)
				} else {
					t.Errorf("actual %v doesn't equal to expected %v", name, "")
				}
			}
			// check that map has been drained
			if !reflect.DeepEqual(tcase.r, make(map[string][]byte)) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.r, make(map[string][]byte))
			}
		})
	}
}
