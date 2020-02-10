package fmts

import (
	"1pkg/gopium"
	"go/ast"
	"go/token"
)

// StructAST defines abstraction for
// formatting gopium.Struct to *ast.TypeSpec
type StructAST func(gopium.Struct) (*ast.TypeSpec, error)

// PrettyJson defines json.Marshal
// with json.Indent TypeFormat implementation
func SimpleASTTypeSpec(st gopium.Struct) (*ast.TypeSpec, error) {
	node := &ast.TypeSpec{}
	node.Name = &ast.Ident{
		Name: st.Name,
		Obj: &ast.Object{
			Kind: ast.Typ,
			Name: st.Name,
		},
	}
	tp := &ast.StructType{
		Fields: &ast.FieldList{
			List: make([]*ast.Field, len(st.Fields)),
		},
	}
	for _, f := range st.Fields {
		astf := &ast.Field{}
		astf.Names = []*ast.Ident{
			&ast.Ident{
				Name: f.Name,
				Obj: &ast.Object{
					Kind: ast.Var,
					Name: f.Name,
				},
			},
		}
		astf.Tag = &ast.BasicLit{
			Kind:  token.STRING,
			Value: f.Tag,
		}
		// TODO type, comment and doc
		tp.Fields.List = append(tp.Fields.List, astf)
	}
	node.Type = tp
	return node, nil
}
