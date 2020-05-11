package walkers

import (
	"context"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"

	"golang.org/x/sync/errgroup"
)

// list of wout presets
var (
	jsonstd = wout{
		fmt:    fmtio.Jsonb,
		writer: fmtio.Stdout,
	}
	xmlstd = wout{
		fmt:    fmtio.Xmlb,
		writer: fmtio.Stdout,
	}
	csvstd = wout{
		fmt:    fmtio.Csvb(fmtio.Buffer()),
		writer: fmtio.Stdout,
	}
	jsonfiles = wout{
		fmt:    fmtio.Jsonb,
		writer: fmtio.File("json"),
	}
	xmlfiles = wout{
		fmt:    fmtio.Xmlb,
		writer: fmtio.File("xml"),
	}
	csvfiles = wout{
		fmt:    fmtio.Csvb(fmtio.Buffer()),
		writer: fmtio.File("csv"),
	}
)

// wout defines packages walker out implementation
type wout struct {
	// inner visiting parameters
	fmt    fmtio.Bytes
	writer fmtio.Writer
	// external visiting parameters
	parser  gopium.TypeParser
	exposer gopium.Exposer
	deep    bool
	bref    bool
}

// With erich wast walker with external visiting parameters
// parser, exposer instances and additional visiting flags
func (w wout) With(p gopium.TypeParser, exp gopium.Exposer, deep, bref bool) wout {
	w.parser = p
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
			if visited.Err != nil {
				return visited.Err
			}
			// just process with write call
			// in case any error happened
			// just return error back
			return w.write(visited.ID, visited.Loc, visited.R)
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
	// generate writer
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
