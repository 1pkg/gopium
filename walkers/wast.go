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
	"1pkg/gopium/fmtio"
)

// list of wast presets
var (
	fsptnstd = wast{
		apply:   apply.SFN,
		persist: persist.AsyncFiles,
		writer:  fmtio.Stdout,
	}
	fsptngo = wast{
		apply:   apply.SFN,
		persist: persist.AsyncFiles,
		writer:  fmtio.FileGo,
	}
	fsptngopium = wast{
		apply:   apply.SFN,
		persist: persist.AsyncFiles,
		writer:  fmtio.FileGopium,
	}
)

// wast defines packages walker sync ast implementation
// that uses pkgs.Parser to parse packages types data
// astutil to update ast to results from strategy
type wast struct {
	parser  gopium.Parser
	exposer gopium.Exposer
	apply   astutil.Apply
	print   astutil.Print
	persist astutil.Persist
	writer  fmtio.Writer
}

// With erich wast walker with parser, exposer, and ref instance
func (w wast) With(parser gopium.Parser, exposer gopium.Exposer, print astutil.Print) wast {
	w.parser = parser
	w.exposer = exposer
	w.print = print
	return w
}

// Visit wast implementation
// uses visit function helper
// to go through all structs decls inside the package
// and apply strategy to them to get results
// then overrides os.Files with updated ast
// builded from strategy results
func (w wast) Visit(
	ctx context.Context,
	regex *regexp.Regexp,
	stg gopium.Strategy,
	deep, backref bool,
) error {
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
	gvisit := visit(
		regex,
		stg,
		w.exposer,
		loc,
		ch,
		deep,
		backref,
	)
	// run visiting in separate goroutine
	go gvisit(ctx, pkg.Scope())
	// prepare struct storage
	h := make(collections.Hierarchic)
	for applied := range ch {
		// in case any error happened just return error
		// it cancels context automatically
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

// write wast helps apply
// sync and persist to format strategy results
// by updating os.Files
func (w wast) write(ctx context.Context, h collections.Hierarchic) error {
	// use parser to parse ast pkg data
	pkg, loc, err := w.parser.ParseAst(ctx)
	if err != nil {
		return err
	}
	// run ast apply with strategy result
	// to update ast.Package on the parsed ast.Package
	// in case any error happened just return error back
	pkg, err = w.apply(ctx, pkg, loc, h)
	if err != nil {
		return err
	}
	// run persist helper
	return w.persist(ctx, w.writer, w.print, pkg, loc)
}
