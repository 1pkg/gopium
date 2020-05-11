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
	JsonbStd  gopium.WalkerName = "json_std"
	XmlbStd   gopium.WalkerName = "xml_std"
	CsvbStd   gopium.WalkerName = "csv_std"
	JsonbFile gopium.WalkerName = "json_file"
	XmlbFile  gopium.WalkerName = "xml_file"
	CsvbFile  gopium.WalkerName = "csv_file"
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
	case JsonbFile:
		return jsonfiles.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case XmlbFile:
		return xmlfiles.With(
			b.Parser,
			b.Exposer,
			b.Deep,
			b.Bref,
		), nil
	case CsvbFile:
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
