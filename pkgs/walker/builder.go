package walker

import (
	"fmt"

	"1pkg/gopium"
	"1pkg/gopium/pkgs"
	"1pkg/gopium/pkgs/read"
	"1pkg/gopium/pkgs/write"
)

// Builder defines types gopium.WalkerBuilder implementation
// that uses pkgs.Parser as an parser and other builder
type Builder struct {
	rb read.Builder
	wb write.Builder
}

// NewBuilder creates instance of Builder
// and requires pkgs.Parser to pass it to other builde
func NewBuilder(parser pkgs.Parser) Builder {
	return Builder{
		rb: read.NewBuilder(parser),
		wb: write.NewBuilder(parser),
	}
}

// Build Builder implementation
func (b Builder) Build(name gopium.WalkerName) (gopium.Walker, error) {
	w, err := b.rb.Build(name)
	if err == nil {
		return w, nil
	}
	w, err = b.wb.Build(name)
	if err == nil {
		return w, nil
	}
	return nil, fmt.Errorf("walker %q wasn't found", name)
}
