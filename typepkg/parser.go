package typepkg

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/1pkg/gopium/gopium"

	"golang.org/x/tools/go/packages"
)

// ParserXToolPackagesAst defines
// gopium parser default implementation
// that uses "golang.org/x/tools/go/packages"
// to collect package types
// and "go/parser" to collect ast package
//
// Note: ParserXToolPackagesAst is big struct
// so it should be passed via pointer
type ParserXToolPackagesAst struct {
	Pattern    string            `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Path       string            `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Root       string            `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	BuildEnv   []string          `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	BuildFlags []string          `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	ModeTypes  packages.LoadMode `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	ModeAst    parser.Mode       `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_          [16]byte          `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 128 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; struct ptr scan size: 80 bytes; - 🌺 gopium @1pkg

// ParseTypes ParserXToolPackagesAst implementation
func (p *ParserXToolPackagesAst) ParseTypes(ctx context.Context, _ ...byte) (*types.Package, gopium.Locator, error) {
	// manage context actions
	// in case of cancelation
	// stop parse and return error back
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}
	// create packages.Config obj
	fset := token.NewFileSet()
	dir := filepath.Join(p.Root, p.Path)
	cfg := &packages.Config{
		Fset:       fset,
		Context:    ctx,
		Dir:        dir,
		Mode:       p.ModeTypes,
		Env:        p.BuildEnv,
		BuildFlags: p.BuildFlags,
		Tests:      true,
	}
	// use packages.Load
	pkgs, err := packages.Load(cfg, "")
	// on any error just propagate it
	if err != nil {
		return nil, nil, err
	}
	// prepare relative path for package
	// by splitting path by src directory
	path := p.Path
	src := fmt.Sprintf("src%s", string(os.PathSeparator))
	if parts := strings.SplitN(p.Path, src, 2); len(parts) == 2 {
		path = parts[1]
		// fix package name for windows
		path = strings.ReplaceAll(path, string(os.PathSeparator), "/")
	}
	// check parse results
	// note: len of pkgs should be equal to
	// - either 1 (pkg contains no tests)
	// - or 3 (pkg contains tests)
	// see go packages config test description
	if plen := len(pkgs); plen >= 1 && validPackage(pkgs[0].String(), p.Pattern, path) {
		switch plen {
		case 1:
			return pkgs[0].Types, NewLocator(fset), nil
		default:
			return pkgs[1].Types, NewLocator(fset), nil
		}
	}
	return nil, nil, fmt.Errorf("types package %q wasn't found at %q", p.Pattern, dir)
}

// ParseAst ParserXToolPackagesAst implementation
func (p *ParserXToolPackagesAst) ParseAst(ctx context.Context, src ...byte) (*ast.Package, gopium.Locator, error) {
	// manage context actions
	// in case of cancelation
	// stop parse and return error back
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}
	fset := token.NewFileSet()
	// if src was provided
	// use parser parse file
	// in memory and return
	// artificial package
	if len(src) > 0 {
		file, err := parser.ParseFile(
			fset,
			"",
			string(src),
			p.ModeAst,
		)
		// on any error just propagate it
		if err != nil {
			return nil, nil, err
		}
		return &ast.Package{
			Name: "pkg",
			Files: map[string]*ast.File{
				"file": file,
			},
		}, NewLocator(fset), err
	}
	// otherwise use parser parse dir
	dir := filepath.Join(p.Root, p.Path)
	pkgs, err := parser.ParseDir(
		fset,
		dir,
		nil,
		p.ModeAst,
	)
	// on any error just propagate it
	if err != nil {
		return nil, nil, err
	}
	// check parse results
	// note: ast requires only last element of the path
	path := filepath.Base(p.Path)
	for pckgn, pckg := range pkgs {
		if validPackage(pckgn, p.Pattern, path) {
			return pckg, NewLocator(fset), nil
		}
	}
	return nil, nil, fmt.Errorf("ast package %q wasn't found at %q", p.Pattern, dir)
}

// validPackage checks package name against provided package pattern and path.
// It preformns shallow sanity package name check to be able to validate
// packages outside of gopath and versioned packages.
func validPackage(pckg string, pattern string, path string) bool {
	// to be extra sure take last component of the path.
	path = filepath.Base(path)
	// first check package name directly
	if direct := strings.HasSuffix(pckg, pattern) ||
		strings.HasSuffix(pckg, path); direct {
		return true
	}
	// then filter out package version
	parts := strings.Split(pckg, "/")
	fparts := make([]string, 0, len(parts))
	for _, part := range parts {
		var version int64
		if _, err := fmt.Sscanf(part, "v%d", &version); err == nil {
			continue
		}
		fparts = append(fparts, part)
	}
	pckg = strings.Join(fparts, "/")
	// and try to recheck again
	return strings.HasSuffix(pckg, pattern) ||
		strings.HasSuffix(pckg, path)
}
