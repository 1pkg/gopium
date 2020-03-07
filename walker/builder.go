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
// that uses pkgs.Parser as an parser and related walkers
type Builder struct {
	parser gopium.Parser
}

// NewBuilder creates instance of Builder
// and requires pkgs.Parser to pass it to related walkers
func NewBuilder(parser gopium.Parser) Builder {
	return Builder{parser: parser}
}

// Build Builder implementation
func (b Builder) Build(name gopium.WalkerName) (gopium.Walker, error) {
	switch name {
	case PrettyJsonStd:
		return wout{
			parser: b.parser,
			fmt:    fmts.PrettyJson,
			writer: os.Stdout,
		}, nil
	case UpdateAst:
		return wuast{
			parser: b.parser,
			fmt:    fmts.FSA,
		}, nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
