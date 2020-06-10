package walkers

import (
	"context"
	"regexp"

	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
	"1pkg/gopium/fmtio/astutil"
	"1pkg/gopium/gopium"
)

// list of wast presets
var (
	aststd = wast{
		apply:     astutil.UFFN,
		persister: astutil.Package{},
		writer:    fmtio.Origin{Writter: fmtio.Stdout{}},
	}
	astgo = wast{
		apply:     astutil.UFFN,
		persister: astutil.Package{},
		writer:    fmtio.Origin{Writter: fmtio.Files{Ext: fmtio.GO}},
	}
	astgotree = wast{
		apply:     astutil.UFFN,
		persister: astutil.Package{},
		writer:    &fmtio.Suffix{Writter: fmtio.Files{Ext: fmtio.GO}, Suffix: gopium.NAME},
	}
	astgopium = wast{
		apply:     astutil.UFFN,
		persister: astutil.Package{},
		writer:    fmtio.Origin{Writter: fmtio.Files{Ext: fmtio.GOPIUM}},
	}
)

// wast defines packages walker ast sync implementation
type wast struct {
	// inner visiting parameters
	apply     gopium.Apply
	persister gopium.Persister
	writer    gopium.CategoryWriter
	// external visiting parameters
	parser  gopium.Parser
	exposer gopium.Exposer
	printer gopium.Printer
	deep    bool
	bref    bool
}

// With erich wast walker with external visiting parameters
// parser, exposer, printer instances and additional visiting flags
func (w wast) With(xp gopium.Parser, exp gopium.Exposer, p gopium.Printer, deep bool, bref bool) wast {
	w.parser = xp
	w.exposer = exp
	w.printer = p
	w.deep = deep
	w.bref = bref
	return w
}

// Visit wast implementation uses visit function helper
// to go through all structs decls inside the package
// and applies strategy to them to get results,
// then overrides ast files with astutil helpers
func (w wast) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
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
	gvisit := with(w.exposer, loc, w.bref).
		visit(regex, stg, ch, w.deep)
	// prepare separate cancelation
	// context for visiting
	gctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// run visiting in separate goroutine
	go gvisit(gctx, pkg.Scope())
	// prepare struct storage
	h := collections.NewHierarchic("")
	for applied := range ch {
		// in case any error happened
		// just return error back
		// it auto cancels context
		if applied.Err != nil {
			return applied.Err
		}
		// push struct to storage
		h.Push(applied.ID, applied.Loc, applied.R)
	}
	// run sync write
	// with collected strategies results
	return w.write(gctx, h)
}

// write wast helps to sync and persist
// strategies results to ast files
func (w wast) write(ctx context.Context, h collections.Hierarchic) error {
	// skip empty writes
	if h.Len() == 0 {
		return nil
	}
	// use parser to parse ast pkg data
	pkg, loc, err := w.parser.ParseAst(ctx)
	if err != nil {
		return err
	}
	// run ast apply with strategy result
	// to update ast.Package
	// in case any error happened
	// just return error back
	pkg, err = w.apply(ctx, pkg, loc, h)
	if err != nil {
		return err
	}
	// add writer root category
	// in case any error happened
	// just return error back
	if err := w.writer.Category(h.Rcat()); err != nil {
		return err
	}
	// run persister with printer
	// in case any error happened
	// just return error back
	return w.persister.Persist(ctx, w.printer, w.writer, loc, pkg)
}
