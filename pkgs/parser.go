package pkgs

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"

	"golang.org/x/tools/go/packages"
)

// Parser defines abstraction for
// types packages parsing processor
type TypeParser interface {
	ParseTypes(context.Context) (*types.Package, error)
}

// Parser defines abstraction for
// ast packages parsing processor
type ASTParser interface {
	ParseAST(context.Context) (*ast.Package, error)
}

// Parser defines abstraction for packages parsing processor
type Parser interface {
	TypeParser
	ASTParser
}

// ParserXToolPackagesAST defines packages Parser default implementation
// that uses "golang.org/x/tools/go/packages" packages.Load with cfg to collect package types
// and uses "go/parser" parser.ParseDir to collect ast package
type ParserXToolPackagesAST struct {
	Pattern    string
	AbsDir     string
	ModeTypes  packages.LoadMode
	ModeAST    parser.Mode
	BuildEnv   []string
	BuildFlags []string
}

// ParseTypes ParserXToolPackagesAST implementation
func (p ParserXToolPackagesAST) ParseTypes(ctx context.Context) (*types.Package, error) {
	// create packages.Config obj
	cfg := &packages.Config{
		Fset:       token.NewFileSet(),
		Context:    ctx,
		Dir:        p.AbsDir,
		Mode:       p.ModeTypes,
		Env:        p.BuildEnv,
		BuildFlags: p.BuildFlags,
		Tests:      true,
	}
	// use packages.Load
	pkgs, err := packages.Load(cfg, p.Pattern)
	if err != nil {
		return nil, err
	}
	// check parse results
	if len(pkgs) != 1 || pkgs[0].String() != p.Pattern {
		return nil, fmt.Errorf("package %q wasn't found", p.Pattern)
	}
	return pkgs[0].Types, nil
}

// ParseAST ParserXToolPackagesAST implementation
func (p ParserXToolPackagesAST) ParseAST(ctx context.Context) (*ast.Package, error) {
	// use parser.ParseDir
	pkgs, err := parser.ParseDir(
		token.NewFileSet(),
		path.Join(p.AbsDir, p.Pattern),
		nil,
		p.ModeAST,
	)
	if err != nil {
		return nil, err
	}
	// check parse results
	pkg, ok := pkgs[path.Base(p.Pattern)]
	if !ok {
		return nil, fmt.Errorf("package %q wasn't found", p.Pattern)
	}
	return pkg, nil
}

// ParserMock defines packages Parser mock implementation
type ParserMock struct {
	typePkg *types.Package
	astPkg  *ast.Package
}

// ParseTypes ParserMock implementation
func (p ParserMock) ParseTypes(context.Context) (*types.Package, error) {
	return p.typePkg, nil
}

// ParseAST ParserMock implementation
func (p ParserMock) ParseAST(context.Context) (*ast.Package, error) {
	return p.astPkg, nil
}

// ParserError defines packages Parser error implementation
type ParserError struct {
	err error
}

// ParseTypes ParserError implementation
func (p ParserError) ParseTypes(context.Context) (*types.Package, error) {
	return nil, p.err
}

// ParseAST ParserError implementation
func (p ParserError) ParseAST(context.Context) (*ast.Package, error) {
	return nil, p.err
}
