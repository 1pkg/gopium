package walkers

import (
	"context"
	"errors"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"

	"golang.org/x/sync/errgroup"
)

// list of wout presets
var (
	jsonstd = wout{
		fmt:  fmtio.PrettyJson,
		wgen: fmtio.Stdout,
		tp:   fmtio.Json,
	}
	xmlstd = wout{
		fmt:  fmtio.PrettyXml,
		wgen: fmtio.Stdout,
		tp:   fmtio.Xml,
	}
	csvstd = wout{
		fmt:  fmtio.PrettyCsv,
		wgen: fmtio.Stdout,
		tp:   fmtio.Csv,
	}
	jsontf = wout{
		fmt:  fmtio.PrettyJson,
		wgen: fmtio.TempFile,
		tp:   fmtio.Json,
	}
	xmltf = wout{
		fmt:  fmtio.PrettyXml,
		wgen: fmtio.TempFile,
		tp:   fmtio.Xml,
	}
	csvtf = wout{
		fmt:  fmtio.PrettyCsv,
		wgen: fmtio.TempFile,
		tp:   fmtio.Csv,
	}
)

// wout defines packages walker out implementation
// that uses pkgs.TypeParser to parse packages types data
// fmts.TypeFormat to format strategy result
// and io.Writer to write output
type wout struct {
	parser  gopium.TypeParser
	exposer gopium.Exposer
	fmt     fmtio.StructToBytes
	wgen    fmtio.WriterGen
	tp      string
	backref bool
}

// With erich wout walker with parser, exposer, and ref instance
func (w wout) With(parser gopium.Parser, exposer gopium.Exposer, backref bool) wout {
	w.parser = parser
	w.exposer = exposer
	w.backref = backref
	return w
}

// Visit wout implementation
func (w wout) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy, deep bool) error {
	// uses gopium.Visit and gopium.VisitFunc helpers
	// to go through all structs decls inside the package
	// and apply strategy to them to get results
	// then use fmts.TypeFormat to format strategy results
	// and use io.Writer to write results to output

	// check that formatter wasn't set correctly
	if w.fmt == nil {
		return errors.New("formatter wasn't set")
	}
	// check that gen wasn't set correctly
	if w.wgen == nil {
		return errors.New("writter generator wasn't set")
	}
	// use parser to parse types pkg data
	// we don't care about fset
	pkg, loc, err := w.parser.ParseTypes(ctx)
	if err != nil {
		return err
	}
	// create govisit func
	// using gopium.Visit helper
	// and run it on pkg scope
	ch := make(appliedCh)
	gvisit := visit(
		regex,
		stg,
		w.exposer,
		loc,
		ch,
		deep,
		w.backref,
	)
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
	// run visiting in separate goroutine
	go gvisit(gctx, pkg.Scope())
loop:
	// go through results from visit func
	// and write them to buf concurently
	for applied := range ch {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-gctx.Done():
			break loop
		default:
		}
		// create applied copy
		visited := applied
		// run error group write call
		group.Go(func() error {
			// in case any error happened just return error
			// it will cancel context automatically
			if visited.Error != nil {
				return visited.Error
			}
			// just process with write call
			// in case any error happened just return error
			// it will cancel context automatically
			return w.write(visited.ID, visited.Loc, visited.Result)
		})
	}
	// wait until all writers
	// resolve their jobs and
	return group.Wait()
}

// write wout helps apply
// fmts.TypeFormat to format strategy result
// and use io.Writer to write result to output
// or return error in any other case
func (w wout) write(id, loc string, st gopium.Struct) error {
	// apply formatter
	buf, err := w.fmt(st)
	// in case any error happened
	// in formatter and return error
	if err != nil {
		return err
	}
	// generate relevant writer
	writer, err := w.wgen(id, loc, w.tp)
	if err != nil {
		return err
	}
	// apply writter
	_, err = writer.Write(buf)
	// in case any error happened
	// in writer return error
	return err
}
