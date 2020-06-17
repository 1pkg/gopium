package walkers

import (
	"fmt"

	"github.com/1pkg/gopium/gopium"
)

// list of registered types walkers
const (
	// wast walkers
	AstStd    gopium.WalkerName = "ast_std"
	AstGo     gopium.WalkerName = "ast_go"
	AstGoTree gopium.WalkerName = "ast_go_tree"
	AstGopium gopium.WalkerName = "ast_gopium"
	// wout walkers
	FileJsonb gopium.WalkerName = "file_json"
	FileXmlb  gopium.WalkerName = "file_xml"
	FileCsvb  gopium.WalkerName = "file_csv"
	FileMdt   gopium.WalkerName = "file_md_table"
	// wdiff walkers
	SizeAlignFileMdt gopium.WalkerName = "size_align_file_md_table"
	FieldsFileHtmlt  gopium.WalkerName = "fields_file_html_table"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses parser and exposer to pass it to related walkers
type Builder struct {
	Parser  gopium.Parser  `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Exposer gopium.Exposer `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Printer gopium.Printer `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Deep    bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Bref    bool           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_       [14]byte       `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// Build Builder implementation
func (b Builder) Build(name gopium.WalkerName) (gopium.Walker, error) {
	switch name {
	// wast walkers
	case AstStd:
		return aststd.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	case AstGo:
		return astgo.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	case AstGoTree:
		return astgotree.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	case AstGopium:
		return astgopium.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	// wout walkers
	case FileJsonb:
		return filejson.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case FileXmlb:
		return filexml.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case FileCsvb:
		return filecsv.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case FileMdt:
		return filemdt.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	// wdiff walkers
	case SizeAlignFileMdt:
		return safilemdt.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case FieldsFileHtmlt:
		return ffilehtml.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
