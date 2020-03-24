package walker

import (
	"fmt"
	"os"

	"1pkg/gopium"
	"1pkg/gopium/fmts"
)

// List of registered types gopium.WalkerName
var (
	PrettyJsonStd gopium.WalkerName = "PrettyJsonStd"
	UpdateAst     gopium.WalkerName = "UpdateAst"
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
	case PrettyJsonStd:
		return wout{
			parser:  b.parser,
			exposer: b.exposer,
			fmt:     fmts.PrettyJson,
			writer:  os.Stdout,
			backref: b.backref,
		}, nil
	case UpdateAst:
		return wuast{
			parser:  b.parser,
			exposer: b.exposer,
			fmt:     fmts.FSPA,
			backref: b.backref,
		}, nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
