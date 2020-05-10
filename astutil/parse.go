package astutil

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// TODO
type parse func(*token.FileSet, string, parser.Mode) (*ast.File, error)

// TODO
func goparse(fset *token.FileSet, content string, mode parser.Mode) (*ast.File, error) {
	return parser.ParseFile(fset, "", content, mode)
}
