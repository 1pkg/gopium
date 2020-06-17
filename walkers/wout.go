package walkers

import (
	"context"
	"path/filepath"
	"regexp"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/fmtio"
	"github.com/1pkg/gopium/gopium"
)

// list of wout presets
var (
	jsonfile = wout{
		fmt:    fmtio.Jsonb,
		writer: fmtio.File{Name: gopium.NAME, Ext: fmtio.JSON},
	}
	xmlfile = wout{
		fmt:    fmtio.Xmlb,
		writer: fmtio.File{Name: gopium.NAME, Ext: fmtio.XML},
	}
	csvfile = wout{
		fmt:    fmtio.Csvb(fmtio.Buffer()),
		writer: fmtio.File{Name: gopium.NAME, Ext: fmtio.CSV},
	}
	mdtfile = wout{
		fmt:    fmtio.Mdtb,
		writer: fmtio.File{Name: gopium.NAME, Ext: fmtio.MD},
	}
)

// wout defines packages walker out implementation
type wout struct {
	writer  gopium.Writer     `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	parser  gopium.TypeParser `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	exposer gopium.Exposer    `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	fmt     gopium.Bytes      `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	deep    bool              `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	bref    bool              `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_       [6]byte           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// With erich wast walker with external visiting parameters
// parser, exposer instances and additional visiting flags
func (w wout) With(p gopium.TypeParser, exp gopium.Exposer, deep bool, bref bool) wout {
	w.parser = p
	w.exposer = exp
	w.deep = deep
	w.bref = bref
	return w
}

// Visit wout implementation uses visit function helper
// to go through all structs decls inside the package
// and applies strategy to them to get results,
// then uses bytes formatter to format strategy results
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
	// prepare separate cancelation
	// context for visiting
	gctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// run visiting in separate goroutine
	go gvisit(gctx, pkg.Scope())
	// prepare struct storage
	h := collections.NewHierarchic("")
	for applied := range ch {
		// in case any error happened
		// just return error back
		// it auto cancels context
		if applied.Err != nil {
			return applied.Err
		}
		// push struct to storage
		h.Push(applied.ID, applied.Loc, applied.R)
	}
	// run sync write
	// with collected strategies results
	return w.write(gctx, h)
}

// write wout helps to apply formatter
// to format strategies result and writer
// to write result to output
func (w wout) write(ctx context.Context, h collections.Hierarchic) error {
	// skip empty writes
	f := h.Flat()
	if len(f) == 0 {
		return nil
	}
	// apply formatter
	buf, err := w.fmt(f.Sorted())
	// in case any error happened
	// in formatter return error back
	if err != nil {
		return err
	}
	// generate writer
	loc := filepath.Join(h.Rcat(), "gopium")
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
