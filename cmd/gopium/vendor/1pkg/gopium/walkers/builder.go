package walkers

import (
	"fmt"

	"1pkg/gopium"
)

// list of registered types walkers
const (
	JsonStd   gopium.WalkerName = "json_std"
	XmlStd    gopium.WalkerName = "xml_std"
	CsvStd    gopium.WalkerName = "csv_std"
	JsonFiles gopium.WalkerName = "json_files"
	XmlFiles  gopium.WalkerName = "xml_files"
	CsvFiles  gopium.WalkerName = "csv_files"
	AstStd    gopium.WalkerName = "ast_std"
	AstGo     gopium.WalkerName = "ast_go"
	AstGopium gopium.WalkerName = "ast_gopium"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses parser and exposer to pass it to related walkers
type Builder struct {
	parser  gopium.Parser
	exposer gopium.Exposer
	backref bool
}

// NewBuilder creates instance of Builder
// and requires parser and exposer to pass it to related walkers
func NewBuilder(parser gopium.Parser, exposer gopium.Exposer, backref bool) Builder {
	return Builder{
		parser:  parser,
		exposer: exposer,
		backref: backref,
	}
}

// Build Builder implementation
func (b Builder) Build(name gopium.WalkerName) (gopium.Walker, error) {
	switch name {
	case JsonStd:
		return jsonstd.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case XmlStd:
		return xmlstd.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case CsvStd:
		return csvstd.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case JsonFiles:
		return jsontf.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case XmlFiles:
		return xmltf.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case CsvFiles:
		return csvtf.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case AstStd:
		return fsptnstd.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case AstGo:
		return fsptngo.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	case AstGopium:
		return fsptngopium.With(
			b.parser,
			b.exposer,
			b.backref,
		), nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
