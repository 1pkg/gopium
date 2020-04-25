package walkers

import (
	"context"
	"errors"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/gfmtio/gfmt"
	"1pkg/gopium/gfmtio/gio"

	"golang.org/x/sync/errgroup"
)

// list of wout presets
var (
	jsonstd = wout{
		fmt:    gfmt.PrettyJson,
		writer: gio.Stdout,
	}
	xmlstd = wout{
		fmt:    gfmt.PrettyXml,
		writer: gio.Stdout,
	}
	csvstd = wout{
		fmt:    gfmt.PrettyCsv,
		writer: gio.Stdout,
	}
	jsonfiles = wout{
		fmt:    gfmt.PrettyJson,
		writer: gio.FileJson,
	}
	xmlfiles = wout{
		fmt:    gfmt.PrettyXml,
		writer: gio.FileXml,
	}
	csvfiles = wout{
		fmt:    gfmt.PrettyCsv,
		writer: gio.FileCsv,
	}
)

// wout defines packages walker out implementation
type wout struct {
	// inner visiting parameters
	fmt    gfmt.StructToBytes
	writer gio.Writer
	// external visiting parameters
	parser  gopium.TypeParser
	exposer gopium.Exposer
	deep    bool
	bref    bool
}

// With erich wast walker with external visiting parameters
// parser, exposer instances and additional visiting flags
func (w wout) With(pars gopium.Parser, exp gopium.Exposer, deep bool, bref bool) wout {
	w.parser = pars
	w.exposer = exp
	w.deep = deep
	w.bref = bref
	return w
}

// Visit wout implementation uses visit function helper
// to go through all structs decls inside the package
// and applies strategy to them to get results,
// then uses struct to bytes to format strategy results
// and use writer to write results to output
func (w wout) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	// check that parser wasn't set correctly
	if w.parser == nil {
		return errors.New("parser wasn't set")
	}
	// check that exposer wasn't set correctly
	if w.exposer == nil {
		return errors.New("exposer wasn't set")
	}
	// check that formatter wasn't set correctly
	if w.fmt == nil {
		return errors.New("formatter wasn't set")
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
	// using gopium.Visit helper
	// and run it on pkg scope
	ch := make(appliedCh)
	gvisit := with(w.exposer, loc, w.bref).
		visit(regex, stg, ch, w.deep)
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
	// run visiting in separate goroutine
	go gvisit(gctx, pkg.Scope())
	// go through results from visit func
	// and write them to buf concurently
	for applied := range ch {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-gctx.Done():
			return gctx.Err()
		default:
		}
		// create applied copy
		visited := applied
		// run error group write call
		group.Go(func() error {
			// in case any error happened
			// just return error back
			if visited.Error != nil {
				return visited.Error
			}
			// just process with write call
			// in case any error happened
			// just return error back
			return w.write(visited.ID, visited.Loc, visited.Result)
		})
	}
	// wait until all writers
	// resolve their jobs and
	return group.Wait()
}

// write wout helps to apply struct to bytes
// to format strategy result and writer
// to write result to output
func (w wout) write(id, loc string, st gopium.Struct) error {
	// apply formatter
	buf, err := w.fmt(st)
	// in case any error happened
	// in formatter return error back
	if err != nil {
		return err
	}
	// generate relevant writer
	writer, err := w.writer(id, loc)
	if err != nil {
		return err
	}
	// write results and close writer
	// in case any error happened
	// in writer return error
	if _, err := writer.Write(buf); err != nil {
		return err
	}
	return writer.Close()
}
