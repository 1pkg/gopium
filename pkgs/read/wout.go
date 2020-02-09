package read

import (
	"context"
	"errors"
	"io"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/fmts"
	"1pkg/gopium/pkgs"
)

// wout defines packages Walker out implementation
// that uses pkgs.Parser to parse packages data
// fmts.TypeFormat to format strategy result
// and io.Writer to write output
type wout struct {
	parser pkgs.Parser
	fmt    fmts.TypeFormat
	writer io.Writer
}

// VisitTop wout implementation
func (w wout) VisitTop(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	return w.visit(ctx, regex, stg, false)
}

// VisitDeep wout implementation
func (w wout) VisitDeep(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	return w.visit(ctx, regex, stg, true)
}

// visit wout helper that uses
// gopium.Visit and gopium.VisitFunc helpers
// to go through all structs decls inside the package
// and apply strategy then get result
// then use fmts.TypeFormat to format strategy result
// and in the end use io.Writer to write output
func (w wout) visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy, deep bool) error {
	// check that formatter wasn't set correctly
	if w.fmt == nil {
		return errors.New("formatter wasn't set")
	}
	// check that writter wasn't set correctly
	if w.writer == nil {
		return errors.New("writter wasn't set")
	}
	// use parser to parse types pkg data
	// and skip ast pkg data
	tpkg, _, err := w.parser.Parse(ctx)
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
	visit(nctx, tpkg.Scope())
	// go through results from visit func
	for ster := range ch {
		// in case any error happened in strategy
		// cancel context and return error
		if ster.Error != nil {
			cancel()
			return ster.Error
		}
		// apply formatter
		buf, err := w.fmt(ster.Struct)
		// in case any error happened in formatter
		// cancel context and return error
		if err != nil {
			cancel()
			return err
		}
		// apply writter
		_, err = w.writer.Write(buf)
		// in case any error happened in formatter
		// cancel context and return error
		if err != nil {
			cancel()
			return err
		}
	}
	// we can sefely cancel context here
	// as walk is done successfully
	// and return nil error
	cancel()
	return nil
}
