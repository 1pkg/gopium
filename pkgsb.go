package gopium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"os"
)

// StgName defines registred strategy name abstraction
type StgName string

var (
	TypeInfoJsonStdOut StgName = "PkgsTiOut-JsonStd"
)

// StgBuilder defines builder abstraction
// that helps to create strategy by name
type StgBuilder interface {
	Build(StgName) (Strategy, error)
}

// Pkgsb defines package strategy builder implementation
// that uses type info extractor abstraction to build strategies
type Pkgsb TiExt

// Build package strategy builder implementation
func (sb Pkgsb) Build(stgnm StgName) (Strategy, error) {
	var exec func(context.Context, string, *types.Struct, *token.FileSet, TiExt) error
	switch stgnm {
	case TypeInfoJsonStdOut:
		f := func(i interface{}) ([]byte, error) {
			r, err := json.Marshal(i)
			if err != nil {
				return nil, err
			}
			var buf bytes.Buffer
			err = json.Indent(&buf, r, "", "\t")
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
		exec = PkgsTiOut{
			tim: make(PkgsTiMap),
			w:   os.Stdout,
			f:   f,
		}.Execute
	default:
		return nil, fmt.Errorf("strategy %q wasn't found", stgnm)
	}

	return func(ctx context.Context, nm string, st *types.Struct, fset *token.FileSet) error {
		return exec(ctx, nm, st, fset, TiExt(sb))
	}, nil
}
