package fmtio

import (
	"bytes"
	"context"
	"go/ast"
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/tests/data"
)

func TestGoprinter(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		p        gopium.Parser
		indent   int
		tabwidth int
		usespace bool
		ctx      context.Context
		r        []byte
		err      error
	}{
		"empty pkg should print nothing": {
			p:        data.NewParser("empty"),
			indent:   0,
			tabwidth: 4,
			usespace: false,
			ctx:      context.Background(),
			r: []byte(`
//+build tests_data

package empty
`),
		},
		"single struct pkg should print the struct": {
			p:        data.NewParser("single"),
			indent:   0,
			tabwidth: 4,
			usespace: false,
			ctx:      context.Background(),
			r: []byte(`
//+build tests_data

package single

type Single struct {
	A	string
	B	string
	C	string
}
`),
		},
		"single struct pkg should print the struct with indent": {
			p:        data.NewParser("single"),
			indent:   1,
			tabwidth: 8,
			usespace: true,
			ctx:      context.Background(),
			r: []byte(`
        //+build tests_data

        package single

        type Single struct {
                A       string
                B       string
                C       string
        }
`),
		},
		"single struct pkg should print nothing on canceled context": {
			p:        data.NewParser("single"),
			indent:   0,
			tabwidth: 4,
			usespace: false,
			ctx:      cctx,
			r:        []byte{},
			err:      context.Canceled,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			buf := &bytes.Buffer{}
			pkg, loc, err := tcase.p.ParseAst(context.Background())
			if !reflect.DeepEqual(err, nil) {
				t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
			}
			// exec
			p := NewGoprinter(tcase.indent, tcase.tabwidth, tcase.usespace)
			// get the only package file
			var file ast.Node
			for _, file = range pkg.Files {
			}
			err = p.Print(tcase.ctx, buf, loc.Root(), file)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// format actual and expected identically
			actual := strings.Trim(string(buf.Bytes()), "\n")
			expected := strings.Trim(string(tcase.r), "\n")
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("name %v actual %v doesn't equal to expected %v", name, actual, expected)
			}
		})
	}
}
