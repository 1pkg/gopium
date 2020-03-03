package walker

import (
	"context"
	"errors"
	"io"
	"regexp"
	"sync"

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
	fmt    fmts.StructFormat
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
	// we don't care about fset
	pkg, _, err := w.parser.ParseTypes(ctx)
	if err != nil {
		return err
	}
	// create gopium.VisitFunc
	// from gopium.Visit helper
	// and run it on pkg scope
	ch := make(gopium.VisitedStructCh)
	visit := gopium.Visit(regex, stg, ch, deep)
	// create separate cancelation context for visiting
	// and defer cancelation func
	nctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// we also need to have separate
	// wait group and error chan
	// to sync concurent writes
	var wg sync.WaitGroup
	wch := make(chan error)
	// run visiting in separate goroutine
	go visit(nctx, pkg.Scope())
loop:
	// go through results from visit func
	// and write them to buf concurently
	for sterr := range ch {
		// in case any error happened just return error
		// it will cancel context automatically
		if sterr.Error != nil {
			return sterr.Error
		}
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
		go func(st gopium.Struct) {
			// decrement writers counter
			defer wg.Done()
			// execute write subaction
			err := w.write(st)
			// in case any error happened put error to chan
			// and cancel context immediately
			if err != nil {
				wch <- err
				cancel()
				return
			}
		}(sterr.Result)
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

// visit wout helps apply
// fmts.TypeFormat to format strategy result
// and use io.Writer to write result to output
// or return error in any other case
func (w wout) write(st gopium.Struct) error {
	// apply formatter
	buf, err := w.fmt(st)
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
