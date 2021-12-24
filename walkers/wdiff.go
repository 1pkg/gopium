package walkers

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/fmtio"
	"github.com/1pkg/gopium/gopium"
)

// list of wdiff presets
var (
	safilemdt = wdiff{
		fmt:    fmtio.SizeAlignMdt,
		writer: fmtio.File{Name: gopium.NAME, Ext: fmtio.MD},
	}
	ffilehtml = wdiff{
		fmt:    fmtio.FieldsHtmlt,
		writer: fmtio.File{Name: gopium.NAME, Ext: fmtio.HTML},
	}
)

// wdiff defines packages walker difference implementation
type wdiff struct {
	writer  gopium.Writer     `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	parser  gopium.TypeParser `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	exposer gopium.Exposer    `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	fmt     gopium.Diff       `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	deep    bool              `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	bref    bool              `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_       [6]byte           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; struct ptr scan size: 56 bytes; - 🌺 gopium @1pkg

// With erich wast walker with external visiting parameters
// parser, exposer instances and additional visiting flags
func (w wdiff) With(p gopium.TypeParser, exp gopium.Exposer, deep bool, bref bool) wdiff {
	w.parser = p
	w.exposer = exp
	w.deep = deep
	w.bref = bref
	return w
}

// Visit wdiff implementation uses visit function helper
// to go through all structs decls inside the package
// and applies strategy to them to get results,
// then uses diff formatter to format strategy results
// and use writer to write results to output
func (w wdiff) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
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
	// prepare separate cancelation
	// context for visiting
	gctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// run visiting in separate goroutine
	go gvisit(gctx, pkg.Scope())
	// prepare struct storages
	ho, hr := collections.NewHierarchic(""), collections.NewHierarchic("")
	for applied := range ch {
		// in case any error happened
		// just return error back
		// it auto cancels context
		if applied.Err != nil {
			return applied.Err
		}
		// push structs to storages
		ho.Push(applied.ID, applied.Loc, applied.O)
		hr.Push(applied.ID, applied.Loc, applied.R)
	}
	// run sync write
	// with collected results
	return w.write(gctx, ho, hr)
}

// write wast helps to apply formatter
// to format strategies results and writer
// to write result to output
func (w wdiff) write(ctx context.Context, ho collections.Hierarchic, hr collections.Hierarchic) error {
	// skip empty writes
	if ho.Len() == 0 || hr.Len() == 0 {
		return nil
	}
	// apply formatter
	buf, err := w.fmt(ho, hr)
	// in case any error happened
	// in formatter return error back
	if err != nil {
		return err
	}
	// generate writer
	loc := filepath.Join(ho.Rcat(), "gopium")
	writer, err := w.writer.Generate(loc)
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
