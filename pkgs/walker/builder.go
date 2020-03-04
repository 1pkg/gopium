package walker

import (
	"fmt"
	"os"

	"1pkg/gopium"
	"1pkg/gopium/fmts"
	"1pkg/gopium/pkgs"
)

// List of registred types gopium.WalkerName
var (
	WalkerOutPrettyJsonStd gopium.WalkerName = "WalkerOut-PrettyJsonStd"
	WalkerUpdateAst        gopium.WalkerName = "WalkerUpdateAst"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses pkgs.Parser as an parser and related walkers
type Builder struct {
	parser pkgs.Parser
}

// NewBuilder creates instance of Builder
// and requires pkgs.Parser to pass it to related walkers
func NewBuilder(parser pkgs.Parser) Builder {
	return Builder{parser: parser}
}

// Build Builder implementation
func (b Builder) Build(name gopium.WalkerName) (gopium.Walker, error) {
	switch name {
	case WalkerOutPrettyJsonStd:
		return wout{
			parser: b.parser,
			fmt:    fmts.PrettyJson,
			writer: os.Stdout,
		}, nil
	case WalkerUpdateAst:
		return wuast(b), nil
	default:
		return nil, fmt.Errorf("walker %q wasn't found", name)
	}
}
