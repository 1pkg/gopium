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
		fmt:    fmtio.PrettyJson,
		writer: fmtio.Stdout,
	}
	xmlstd = wout{
		fmt:    fmtio.PrettyXml,
		writer: fmtio.Stdout,
	}
	csvstd = wout{
		fmt:    fmtio.PrettyCsv,
		writer: fmtio.Stdout,
	}
	jsontf = wout{
		fmt:    fmtio.PrettyJson,
		writer: fmtio.FileJson,
	}
	xmltf = wout{
		fmt:    fmtio.PrettyXml,
		writer: fmtio.FileXml,
	}
	csvtf = wout{
		fmt:    fmtio.PrettyCsv,
		writer: fmtio.FileCsv,
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
	writer  fmtio.Writer
}

// With erich wout walker with parser, exposer, and ref instance
func (w wout) With(parser gopium.Parser, exposer gopium.Exposer) wout {
	w.parser = parser
	w.exposer = exposer
	return w
}

// Visit wout implementation
// uses visit function helper
// to go through all structs decls inside the package
// and apply strategy to them to get results
// then use fmts.TypeFormat to format strategy results
// and use io.Writer to write results to output
func (w wout) Visit(
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
	gvisit := visit(
		regex,
		stg,
		w.exposer,
		loc,
		ch,
		deep,
		backref,
	)
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
			// in case any error happened just return error
			// it cancels context automatically
			if visited.Error != nil {
				return visited.Error
			}
			// just process with write call
			// in case any error happened just return error
			// it cancels context automatically
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
	writer, err := w.writer(id, loc)
	if err != nil {
		return err
	}
	// apply writer
	_, err = writer.Write(buf)
	// in case any error happened
	// in writer return error
	return err
}
