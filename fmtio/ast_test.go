package fmtio

import (
	"bytes"
	"errors"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/astutil/print"
	"1pkg/gopium/tests/mocks"
)

func TestAst(t *testing.T) {
	// prepare
	p := print.GoPrinter(0, 4, false)
	table := map[string]struct {
		f   Ast
		ts  *ast.TypeSpec
		st  gopium.Struct
		r   []byte
		err error
	}{
		"not struct type should lead to error": {
			f: FSPT,
			ts: &ast.TypeSpec{
				Name: &ast.Ident{
					Name: "test",
				},
				Type: &ast.ArrayType{
					Len: &ast.BasicLit{
						Kind:  token.INT,
						Value: "10",
					},
					Elt: &ast.Ident{
						Name: "string",
					},
				},
			},
			err: errors.New(`type "test" is not valid structure`),
		},
		"error from ast func should be transfered": {
			f: combine(flatten, mocks.Ast{Err: errors.New("test")}.Ast),
			ts: &ast.TypeSpec{
				Name: &ast.Ident{
					Name: "test",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{},
					},
				},
			},
			err: errors.New("test"),
		},
		"non flat struct with tag should be flatten correctly": {
			f: FSPT,
			ts: &ast.TypeSpec{
				Name: &ast.Ident{
					Name: "test",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "test-1",
									},
									{
										Name: "test-2",
									},
									{
										Name: "test-3",
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
										Name: "test-4",
									},
									{
										Name: "test-5",
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
								Type: &ast.Ident{
									Name: "float32",
								},
								Tag: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "embedded",
								},
							},
						},
					},
				},
			},
			st: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test-1",
						Type: "string",
						Tag:  "test-1",
					},
					{
						Name: "test-2",
						Type: "string",
						Tag:  "test-2",
					},
					{
						Name: "test-3",
						Type: "string",
						Tag:  "test-3",
					},
					{
						Name: "test-4",
						Type: "int64",
					},
					{
						Name: "test-5",
						Type: "int64",
					},
					{
						Type:     "float32",
						Embedded: true,
					},
				},
			},
			r: []byte(`
test struct {
	test-1	string	'test-1'
	test-2	string	'test-2'
	test-3	string	'test-3'
	test-4	int64
	test-5	int64// random
	// random
	float32	embedded
}
`),
		},
		"struct with excess paddings and fields should be filtered and sorted": {
			f: FSPT,
			ts: &ast.TypeSpec{
				Name: &ast.Ident{
					Name: "test",
				},
				Type: &ast.StructType{
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
				},
			},
			st: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test-2",
						Type: "int64",
					},
					{
						Name: "test-1",
						Type: "int64",
					},
				},
			},
			r: []byte(`
test struct {// random
	test-2	int64
	test-1	int64// random
}
`),
		},
		"struct paddings should be synchronized": {
			f: FSPT,
			ts: &ast.TypeSpec{
				Name: &ast.Ident{
					Name: "test",
				},
				Type: &ast.StructType{
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
				},
			},
			st: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test-1",
						Type: "int64",
					},
					{
						Name: "test-2",
						Type: "int64",
					},
					{
						Name: "_",
						Type: "[10]byte",
						Size: 10,
						Tag:  "btag",
					},
					{
						Name: "_",
						Type: "float64",
						Size: 8,
					},
				},
			},
			r: []byte(`
test struct {// random
	test-1	int64
	test-2	int64
	_		[// random
	10]byte	'btag'
	_	[8]byte
}
`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			var buf bytes.Buffer
			err := tcase.f(tcase.ts, tcase.st)
			perr := p(&buf, token.NewFileSet(), tcase.ts)
			// check
			if !reflect.DeepEqual(perr, nil) {
				t.Errorf("actual %v doesn't equal to expected %v", perr, nil)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// format actual and expected identically
			actual := strings.Trim(string(buf.Bytes()), "\n")
			expected := strings.ReplaceAll(strings.Trim(string(tcase.r), "\n"), "'", "`")
			if err == nil && !reflect.DeepEqual(actual, expected) {
				t.Errorf("actual %v doesn't equal to expected %v", actual, expected)
			}
		})
	}
}
