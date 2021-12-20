package fmtio

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"testing"
)

func TestGoprinter(t *testing.T) {
	// prepare
	node := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "test-removed",
						},
						{
							Name: "_",
						},
					},
					Type: &ast.Ident{
						Name: "string",
					},
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: "test",
					},
				},
				{
					Names: []*ast.Ident{
						{
							Name: "test-1",
						},
						{
							Name: "test-2",
						},
					},
					Type: &ast.Ident{
						Name: "int64",
					},
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{
							{
								Text: "// random",
							},
						},
					},
				},
				{
					Names: []*ast.Ident{
						{
							Name: "_",
						},
					},
					Type: &ast.Ident{
						Name: "float32",
					},
					Tag: &ast.BasicLit{
						Kind:  token.STRING,
						Value: "tag",
					},
				},
			},
		},
	}
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		indent   int
		tabwidth int
		usespace bool
		ctx      context.Context
		r        []byte
		err      error
	}{
		"single struct pkg should print the struct": {
			indent:   0,
			tabwidth: 4,
			usespace: false,
			ctx:      context.Background(),
			r: []byte(`
struct {
	test-removed, _	string	test// random
	test-1, test-2	int64
	_				float32	tag
}
`),
		},
		"single struct pkg should print the struct with indent": {
			indent:   1,
			tabwidth: 8,
			usespace: true,
			ctx:      context.Background(),
			r: []byte(`
        struct {
                test-removed, _ string  test// random
                test-1, test-2  int64
                _               float32 tag
        }
`),
		},
		"single struct pkg should print nothing on canceled context": {
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
			// exec
			var buf bytes.Buffer
			p := NewGoprinter(tcase.indent, tcase.tabwidth, tcase.usespace)
			err := p.Print(tcase.ctx, &buf, token.NewFileSet(), node)
			// check
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// format actual and expected identically
			actual := strings.Trim(buf.String(), "\n")
			expected := strings.Trim(string(tcase.r), "\n")
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("name %v actual %v doesn't equal to expected %v", name, actual, expected)
			}
		})
	}
}
