package walker

import (
	"context"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"regexp"
	"sync"

	"golang.org/x/tools/go/ast/astutil"

	"1pkg/gopium"
	"1pkg/gopium/pkgs"
)

// wast defines packages Walker AST implementation
// that uses pkgs.Parser to parse packages types data
// astutil to update AST to results from strategy
type wast struct {
	parser pkgs.Parser
}

// VisitTop wast implementation
func (w wast) VisitTop(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	return w.visit(ctx, regex, stg, false)
}

// VisitDeep wast implementation
func (w wast) VisitDeep(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	return w.visit(ctx, regex, stg, true)
}

// visit wast helps with visiting and uses
// gopium.Visit and gopium.VisitFunc helpers
// to go through all structs decls inside the package
// and apply strategy to them to get results
// then overrides os.File list with updated AST
// builded from strategy results
func (w wast) visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy, deep bool) error {
	// use parser to parse types pkg data
	// we don't care about fset
	pkg, _, err := w.parser.ParseTypes(ctx)
	if err != nil {
		return err
	}
	// create gopium.VisitFunc
	// from gopium.Visit helper
	// and run it on pkg scope
	ch := make(gopium.VisitedStructCh)
	// TODO use real gopium.IDFunc impl
	visit := gopium.Visit(regex, stg, nil, ch, deep)
	// create separate cancelation context for visiting
	nctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// run visiting in separate goroutine
	go visit(nctx, pkg.Scope())
	// go through results from visit func
	// we can use concurent writitng too
	// but it's probably redundant
	// as it requires additional level of sync
	// and intense error handling
	var structs []gopium.Struct
	for applied := range ch {
		// in case any error happened just return error
		// it will cancel context automatically
		if applied.Error != nil {
			return applied.Error
		}
		// otherwise collect result
		structs = append(structs, applied.Result)
	}
	// run sync write
	// with collected strategies results
	return w.write(nctx, structs)
}

// write wast helps apply
// updateAST/writeAST to format strategy results
// updating os.File list ASTs
func (w wast) write(ctx context.Context, structs []gopium.Struct) error {
	// use parser to parse ast pkg data
	pkg, fset, err := w.parser.ParseAST(ctx)
	if err != nil {
		return err
	}
	// go through results from visit func
	// we can use concurent updating too
	// but it's probably redundant
	// as it requires additional level of sync
	// and intense error handling
	for _, st := range structs {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		// run updateAST with strategy result
		// on the parsed ast.Package
		// in case any error happened just return error
		pkg, err = w.updateAST(ctx, pkg, st)
		if err != nil {
			return err
		}
	}
	// run async writeAST helper
	return w.writeAST(ctx, pkg, fset)
}

// updateAST helps to update ast.Package
// accordingly to Strategy gopium.Struct result
// synchronously or return error otherwise
func (w wast) updateAST(ctx context.Context, pkg *ast.Package, st gopium.Struct) (*ast.Package, error) {
	// apply astutil.Apply to parsed ast.Package
	// and update structure in AST
	unode := astutil.Apply(pkg, func(c *astutil.Cursor) bool {
		node := c.Node()
		if gendecl, ok := node.(*ast.GenDecl); ok {
			for _, spec := range gendecl.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := ts.Type.(*ast.StructType); ok {
						// TODO use fmts.StructAST to generate AST for gopium.Struct
						return true
					}
				}
			}
		}
		return true

	}, nil)
	// check that updated type is correct
	if upkg, ok := unode.(*ast.Package); ok {
		return upkg, nil
	}
	// in case updated type isn't expected
	return nil, fmt.Errorf("can't update AST for structure %+v", st)
}

// writeAST helps to update os.File list
// accordingly to updated ast.Package
// concurently or return error otherwise
func (w wast) writeAST(ctx context.Context, pkg *ast.Package, fset *token.FileSet) error {
	// create separate cancelation context for writing
	// and defer cancelation func
	nctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// we also need to have separate
	// wait group and error chan
	// to sync concurent writes
	var wg sync.WaitGroup
	wch := make(chan error)
loop:
	// go through all files in package
	// and update them to concurently
	for name, file := range pkg.Files {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-nctx.Done():
			break loop
		default:
		}
		// increment writers counter
		wg.Add(1)
		// concurently update each ast.File to os.File
		go func(fname string, node *ast.File) {
			// decrement writers counter
			defer wg.Done()
			// open os.File for related ast.File
			file, err := os.Create(fname)
			// in case any error happened put error to chan
			// and cancel context immediately
			if err != nil {
				wch <- err
				cancel()
				return
			}
			// write updated ast.File to related os.File
			// use original toke.FileSet to keep format
			err = printer.Fprint(
				file,
				fset,
				node,
			)
			// in case any error happened put error to chan
			// and cancel context immediately
			if err != nil {
				wch <- err
				cancel()
				return
			}
		}(name, file)
	}
	// will wait until all writers
	// resolve their jobs and
	// close error wait ch after
	go func() {
		wg.Wait()
		close(wch)
	}()
	return <-wch
}
