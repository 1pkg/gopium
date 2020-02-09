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
// that uses pkgs.TypeParser to parse packages types data
// fmts.TypeFormat to format strategy result
// and io.Writer to write output
type wout struct {
	parser pkgs.TypeParser
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

// visit wout helps with visiting and uses
// gopium.Visit and gopium.VisitFunc helpers
// to go through all structs decls inside the package
// and apply strategy to them to get results
// then use fmts.TypeFormat to format strategy results
// and use io.Writer to write results to output
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
	tpkg, err := w.parser.ParseTypes(ctx)
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
	go visit(nctx, tpkg.Scope())
	// go through results from visit func
	// we can use concurent writitng too
	// but it's probably redundant
	// as it requires additional level of sync
	// and error handling
	for sterr := range ch {
		err := w.write(ctx, sterr)
		// in case any error happened in writting
		// cancel context and return error
		if err != nil {
			cancel()
			return err
		}
	}
	// we can safely cancel context here
	// as walk is already done successfully
	// and returned nil error
	cancel()
	return nil
}

// visit wout helps apply
// fmts.TypeFormat to format strategy result
// and use io.Writer to write result to output
// or return error in any other case
func (w wout) write(ctx context.Context, sterr gopium.StructError) error {
	// in case any error happened
	// in strategy return error
	if sterr.Error != nil {
		return sterr.Error
	}
	// apply formatter
	buf, err := w.fmt(sterr.Struct)
	// in case any error happened
	// in formatter and return error
	if err != nil {
		return err
	}
	// apply writter
	_, err = w.writer.Write(buf)
	// in case any error happened
	// in writer return error
	return err
}
