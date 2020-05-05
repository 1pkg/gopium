package walkers

import (
	"fmt"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
)

// list of registered types walkers
const (
	// wast walkers
	AstStd    gopium.WalkerName = "ast_std"
	AstGo     gopium.WalkerName = "ast_go"
	AstGopium gopium.WalkerName = "ast_gopium"
	// wout walkers
	JsonStd   gopium.WalkerName = "json_std"
	XmlStd    gopium.WalkerName = "xml_std"
	CsvStd    gopium.WalkerName = "csv_std"
	JsonFiles gopium.WalkerName = "json_files"
	XmlFiles  gopium.WalkerName = "xml_files"
	CsvFiles  gopium.WalkerName = "csv_files"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses parser and exposer to pass it to related walkers
type Builder struct {
	Parser  gopium.Parser
	Exposer gopium.Exposer
	Printer fmtio.Printer
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
	case AstGopium:
		return astgopium.With(
			b.Parser,
			b.Exposer,
			b.Printer,
			b.Deep,
			b.Bref,
		), nil
	// wout walkers
	case JsonStd:
		return jsonstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case XmlStd:
		return xmlstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case CsvStd:
		return csvstd.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case JsonFiles:
		return jsonfiles.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case XmlFiles:
		return xmlfiles.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case CsvFiles:
		return csvfiles.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
