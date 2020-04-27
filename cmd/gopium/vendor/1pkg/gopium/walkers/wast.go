package walkers

import (
	"context"
	"errors"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/astutil/apply"
	"1pkg/gopium/astutil/persist"
	"1pkg/gopium/collections"
	"1pkg/gopium/gfmtio/gio"
)

// list of wast presets
var (
	aststd = wast{
		apply:   apply.SFN,
		persist: persist.AsyncFiles,
		writer:  gio.Stdout,
	}
	astgo = wast{
		apply:   apply.SFN,
		persist: persist.AsyncFiles,
		writer:  gio.FileGo,
	}
	astgopium = wast{
		apply:   apply.SFN,
		persist: persist.AsyncFiles,
		writer:  gio.FileGopium,
	}
)

// wast defines packages walker ast sync implementation
type wast struct {
	// inner visiting parameters
	apply   astutil.Apply
	persist astutil.Persist
	writer  gio.Writer
	// external visiting parameters
	parser  gopium.Parser
	exposer gopium.Exposer
	print   astutil.Print
	deep    bool
	bref    bool
}

// With erich wast walker with external visiting parameters
// parser, exposer, printer instances and additional visiting flags
func (w wast) With(pars gopium.Parser, exp gopium.Exposer, pr astutil.Print, deep bool, bref bool) wast {
	w.parser = pars
	w.exposer = exp
	w.print = pr
	w.deep = deep
	w.bref = bref
	return w
}

// Visit wast implementation uses visit function helper
// to go through all structs decls inside the package
// and applies strategy to them to get results,
// then overrides ast files with astutil helpers
func (w wast) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	// check that parser wasn't set correctly
	if w.parser == nil {
		return errors.New("parser wasn't set")
	}
	// check that exposer wasn't set correctly
	if w.exposer == nil {
		return errors.New("exposer wasn't set")
	}
	// check that apply wasn't set correctly
	if w.apply == nil {
		return errors.New("apply wasn't set")
	}
	// check that print wasn't set correctly
	if w.print == nil {
		return errors.New("print wasn't set")
	}
	// check that persist wasn't set correctly
	if w.persist == nil {
		return errors.New("persist wasn't set")
	}
	// check that writer wasn't set correctly
	if w.writer == nil {
		return errors.New("writer wasn't set")
	}
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
	h := make(collections.Hierarchic)
	for applied := range ch {
		// in case any error happened
		// just return error back
		// it auto cancels context
		if applied.Error != nil {
			return applied.Error
		}
		// push struct to storage
		h.Push(applied.ID, applied.Loc, applied.Result)
	}
	// run sync write
	// with collected strategies results
	return w.write(ctx, h)
}

// write wast helps to sync
// and persist strategy results to ast files
func (w wast) write(ctx context.Context, h collections.Hierarchic) error {
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
	// run persist helper
	// in case any error happened
	// just return error back
	return w.persist(ctx, w.writer, w.print, pkg, loc)
}
