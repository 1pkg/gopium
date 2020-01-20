package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"os"
)

// SgName defines registred strategy name abstraction
type SgName string

var (
	TypeInfoJsonStdOut SgName = "PkgsTiOut-JsonStd"
)

// SgBuilder defines builder abstraction
// that helps to create strategy by name
type SgBuilder interface {
	Build(string) (Strategy, error)
}

// Pkgsb defines package strategy builder implementation
// that uses type info extractor abstraction to build strategies
type Pkgsb TiExt

// Build package strategy builder implementation
func (sb Pkgsb) Build(sgn SgName) (Strategy, error) {
	var exec func(context.Context, *types.Struct, *token.FileSet, TiExt) error
	switch sgn {
	case TypeInfoJsonStdOut:
		exec = PkgsTiOut{
			tim: make(PkgsTiMap),
			w:   os.Stdout,
			f:   json.Marshal,
		}.Execute
	default:
		return nil, fmt.Errorf("strategy `%s` wasn't found", sgn)
	}

	return func(ctx context.Context, st *types.Struct, fset *token.FileSet) error {
		return exec(ctx, st, fset, TiExt(sb))
	}, nil
}
