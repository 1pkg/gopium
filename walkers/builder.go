package walkers

import (
	"fmt"

	"1pkg/gopium"
)

// list of registered types walkers
const (
	// wast walkers
	AstStd    gopium.WalkerName = "ast_std"
	AstGo     gopium.WalkerName = "ast_go"
	AstGoTree gopium.WalkerName = "ast_go_tree"
	AstGopium gopium.WalkerName = "ast_gopium"
	// wout walkers
	JsonbStd  gopium.WalkerName = "json_std"
	XmlbStd   gopium.WalkerName = "xml_std"
	CsvbStd   gopium.WalkerName = "csv_std"
	MdtStd    gopium.WalkerName = "md_table_std"
	JsonbFile gopium.WalkerName = "json_file"
	XmlbFile  gopium.WalkerName = "xml_file"
	CsvbFile  gopium.WalkerName = "csv_file"
	MdtFile   gopium.WalkerName = "md_table_file"
	// wdiff walkers
	SizeDiffMdtStd  gopium.WalkerName = "size_diff_md_table_std"
	SizeDiffMdtFile gopium.WalkerName = "size_diff_md_table_file"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses parser and exposer to pass it to related walkers
type Builder struct {
	Parser  gopium.Parser
	Exposer gopium.Exposer
	Printer gopium.Printer
	Deep    bool
	Bref    bool
}

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
	case JsonbStd:
		return jsonstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case XmlbStd:
		return xmlstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case CsvbStd:
		return csvstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case MdtStd:
		return mdtstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case JsonbFile:
		return jsonfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case XmlbFile:
		return xmlfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case CsvbFile:
		return csvfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case MdtFile:
		return mdtfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	// wdiff walkers
	case SizeDiffMdtStd:
		return satmdstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case SizeDiffMdtFile:
		return satmdfile.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
