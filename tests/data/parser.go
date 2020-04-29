package data

import (
	"fmt"
	"go/build"
	"go/parser"

	"1pkg/gopium"
	"1pkg/gopium/typepkg"

	"golang.org/x/tools/go/packages"
)

// NewParser creates parser for single tests data
func NewParser(pkg string) gopium.Parser {
	return typepkg.ParserXToolPackagesAst{
		Path:       fmt.Sprintf("%s/%s", "src/1pkg/gopium/tests/data", pkg),
		Root:       build.Default.GOPATH,
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments | parser.AllErrors,
		BuildFlags: []string{"-tags=tests_data"},
	}
}
