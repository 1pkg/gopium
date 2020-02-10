package write

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
	ch := make(chan gopium.StructError)
	visit := gopium.Visit(regex, stg, ch, deep)
	// create separate cancelation context for visiting
	nctx, cancel := context.WithCancel(ctx)
	// run visiting in separate goroutine
	go visit(nctx, pkg.Scope())
	// go through results from visit func
	// we can use concurent writitng too
	// but it's probably redundant
	// as it requires additional level of sync
	// and error handling
	var sterrs []gopium.StructError
	for sterr := range ch {
		// manage context actions
		// in case of cancelation break from
		// collecting action
		select {
		default:
			// in case any error happened
			// cancel context and return error
			if sterr.Error != nil {
				cancel()
				return sterr.Error
			}
			// collect result
			sterrs = append(sterrs, sterr)
		case <-ctx.Done():
			cancel()
			return nil
		}
	}
	// we can safely cancel context here
	// as walk is already done successfully
	// and returned nil error
	cancel()
	// run sync write
	// with collected strategies results
	return w.write(ctx, sterrs)
}

// write wast helps apply
// updateAST/writeAST to format strategy results
// updating os.File list ASTs
func (w wast) write(ctx context.Context, sterrs []gopium.StructError) error {
	// use parser to parse ast pkg data
	pkg, fset, err := w.parser.ParseAST(ctx)
	if err != nil {
		return err
	}
	// go through results from visit func
	// we can use concurent updating too
	// but it's probably redundant
	// as it requires additional level of sync
	// and error handling
	for _, sterr := range sterrs {
		// manage context actions
		// in case of cancelation break from
		// writting action
		select {
		default:
			// run updateAST with strategy result
			// on the parsed AST pkg
			// in case any error happened just return error
			pkg, err = w.updateAST(ctx, pkg, sterr)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
	// run async writeAST helper
	return w.writeAST(ctx, pkg, fset)
}

// updateAST helps to update ast.Package
// accordingly to Strategy gopium.StructError result
// synchronously or return error otherwise
func (w wast) updateAST(ctx context.Context, pkg *ast.Package, sterr gopium.StructError) (*ast.Package, error) {
	// in case strategy has any error
	// just return error
	if sterr.Error != nil {
		return nil, sterr.Error
	}
	// apply astutil.Apply to parsed ast.Package
	// and update structure in AST
	unode := astutil.Apply(pkg, func(c *astutil.Cursor) bool {
		return true
	}, nil)
	if upkg, ok := unode.(*ast.Package); ok {
		return upkg, nil
	}
	// in case updated type isn't expected
	return nil, fmt.Errorf("can't update AST for structure %+v", sterr.Struct)
}

// writeAST helps to update os.File list
// accordingly to updated ast.Package
// concurently or return error otherwise
func (w wast) writeAST(ctx context.Context, pkg *ast.Package, fset *token.FileSet) error {
	// create separate cancelation context for writting
	//nolint
	nctx, cancel := context.WithCancel(ctx)
	// wait group writing counter
	var wg sync.WaitGroup
	// channel that keeps writting errors
	errch := make(chan error)
	for name, file := range pkg.Files {
		// increment wait group writing counter anyway
		wg.Add(1)
		// concurently update each ast.File to os.File
		go func(ctx context.Context, fname string, node *ast.File) {
			// decrement wait group writing counter anyway
			defer wg.Done()
			// manage context actions
			// in case of cancelation break from
			// writting action
			select {
			default:
				// open os.File for related ast.File
				// in case of any error put it to errch
				// and cancel context
				file, err := os.Create(fname)
				if err != nil {
					errch <- err
					cancel()
					return
				}
				// write updated ast.File to related os.File
				// use original toke.FileSet to keep format
				// in case of any error put it to errch
				// and cancel context
				err = printer.Fprint(
					file,
					fset,
					node,
				)
				if err != nil {
					errch <- err
					cancel()
					return
				}
			case <-ctx.Done():
				cancel()
				return
			}
		}(nctx, name, file)
	}
	// wait until all writtings are done
	// and close writting errors channel
	wg.Wait()
	// get last error from the channel
	var err error
	select {
	case err = <-errch:
	default:
		// we can safely cancel context here
		// as write is already done successfully
		// and returned nil error
		cancel()
	}
	// close error channel and return last error
	close(errch)
	//nolint
	return err
}
