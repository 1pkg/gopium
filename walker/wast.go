package walker

import (
	"context"
	"fmt"
	"go/ast"
	"go/printer"
	"os"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/fmts"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/ast/astutil"
)

// list of wast presets
var (
	fsptn = wast{
		fmt: fmts.FSPTN,
	}
)

// wast defines packages walker sync ast implementation
// that uses pkgs.Parser to parse packages types data
// astutil to update ast to results from strategy
type wast struct {
	parser  gopium.Parser
	exposer gopium.Exposer
	fmt     fmts.StructToAst
	backref bool
}

// VisitTop wast implementation
func (w wast) VisitTop(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	return w.visit(ctx, regex, stg, false)
}

// VisitDeep wast implementation
func (w wast) VisitDeep(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	return w.visit(ctx, regex, stg, true)
}

// With erich wast walker with parser, exposer, and ref instance
func (w wast) With(parser gopium.Parser, exposer gopium.Exposer, backref bool) wast {
	w.parser = parser
	w.exposer = exposer
	w.backref = backref
	return w
}

// visit wast helps with visiting and uses
// gopium.Visit and gopium.VisitFunc helpers
// to go through all structs decls inside the package
// and apply strategy to them to get results
// then overrides os.File list with updated ast
// builded from strategy results
func (w wast) visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy, deep bool) error {
	// use parser to parse types pkg data
	// we don't care about fset
	pkg, loc, err := w.parser.ParseTypes(ctx)
	if err != nil {
		return err
	}
	// create govisit func
	// using visit helper
	// and run it on pkg scope
	ch := make(appliedCh)
	gvisit := visit(regex, stg, w.exposer, loc.Sum, ch, deep, w.backref)
	// create separate cancelation context for visiting
	nctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// run visiting in separate goroutine
	go gvisit(nctx, pkg.Scope())
	// go through results from visit func
	// we can use concurent writitng too
	// but it's probably redundant
	// as it requires additional level of sync
	// and intense error handling
	structs := make(map[string]gopium.Struct)
	for applied := range ch {
		// in case any error happened just return error
		// it will cancel context automatically
		if applied.Error != nil {
			return applied.Error
		}
		// otherwise collect result
		structs[applied.ID] = applied.Result
	}
	// run sync write
	// with collected strategies results
	return w.write(nctx, structs)
}

// write wast helps apply
// sync and persist to format strategy results
// updating os.File ast list
func (w wast) write(ctx context.Context, structs map[string]gopium.Struct) error {
	// use parser to parse ast pkg data
	pkg, loc, err := w.parser.ParseAst(ctx)
	if err != nil {
		return err
	}
	// go through results from visit func
	// we can use concurent updating too
	// but it's probably redundant
	// as it requires additional level of sync
	// and intense error handling
	for id, st := range structs {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		// run sync with strategy result to update ast.Package
		// on the parsed ast.Package
		// in case any error happened just return error
		pkg, err = w.sync(pkg, loc, id, st)
		if err != nil {
			return err
		}
	}
	// run async persist helper
	return w.persist(ctx, pkg, loc)
}

// sync wast helps to update ast.Package
// accordingly to Strategy gopium.Struct result
// synchronously or return error otherwise
func (w wast) sync(pkg *ast.Package, loc *gopium.Locator, id string, st gopium.Struct) (*ast.Package, error) {
	// tracks error inside astutil.Apply
	var err error
	// apply astutil.Apply to parsed ast.Package
	// and update structure in ast
	snode := astutil.Apply(pkg, func(c *astutil.Cursor) bool {
		if gendecl, ok := c.Node().(*ast.GenDecl); ok {
			for _, spec := range gendecl.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := ts.Type.(*ast.StructType); ok {
						// calculate sum for structure
						// and skip all irrelevant structs
						sum := loc.Sum(ts.Pos())
						if id == sum {
							// apply format to ast
							err = w.fmt(ts, st)
							// in case we have error
							// break iteration
							return err != nil
						}
					}
				}
			}
		}
		return true
	}, nil)
	// in case we had error in astutil.Apply
	// just return it back
	if err != nil {
		return nil, err
	}
	// check that updated type is correct
	if spkg, ok := snode.(*ast.Package); ok {
		return spkg, nil
	}
	// in case updated type isn't expected
	return nil, fmt.Errorf("can't update ast for structure %q", st.Name)
}

// persist wast helps to update os.File list
// accordingly to updated ast.Package
// concurently or return error otherwise
func (w wast) persist(ctx context.Context, pkg *ast.Package, loc *gopium.Locator) error {
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
loop:
	// go through all files in package
	// and update them to concurently
	for fname, file := range pkg.Files {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-gctx.Done():
			break loop
		default:
		}
		// create fname and file copies
		name := fname
		node := file
		// run error group write call
		group.Go(func() error {
			// open os.File for related ast.File
			file, err := os.Create(name)
			// in case any error happened just return error
			// it will cancel context automatically
			if err != nil {
				return err
			}
			// write updated ast.File to related os.File
			// use original toke.FileSet to keep format
			// in case any error happened just return error
			// it will cancel context automatically
			return printer.Fprint(
				file,
				loc.Fset(),
				node,
			)
		})
	}
	// wait until all writers
	// resolve their jobs and
	return group.Wait()
}
