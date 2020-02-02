package strategy

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"io"

	"1pkg/gopium/fmts"
)

// typeOut defines strategy type out implementation
// that uses strategy typeMap result
// format it with custom gopium.Formatter function
// and finaly outputs it to io.Writer
type typeOut struct {
	typeMap
	f fmts.TypeFormat
	w io.Writer
}

// Execute strategy type out implementation
func (to *typeOut) Execute(ctx context.Context, nm string, st *types.Struct, fset *token.FileSet) error {
	// check that formatter and writter were set properly
	if to.f == nil {
		return errors.New("formatter method wasn't set")
	}
	if to.w == nil {
		return errors.New("writter method wasn't set")
	}
	// execute underlying typeMap
	err := to.typeMap.Execute(ctx, nm, st, fset)
	if err != nil {
		return err
	}
	// add struct name to fields typeinfo map
	r := make(map[string]interface{})
	r["Name"] = nm
	r["Fields"] = to.r
	// apply formatter
	buf, err := to.f(r)
	if err != nil {
		return err
	}
	// apply writter
	_, err = to.w.Write(buf)
	return err
}
