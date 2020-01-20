package main

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"io"
)

// Strategy defines action abstraction
// that applies some strategy on types struct
type Strategy func(context.Context, *types.Struct, *token.FileSet) error

// PkgsTiMap defines strategy implementation
// that goes through structure fields
// extracts type info for each field and put it to the map
type PkgsTiMap map[string]TypeInfo

// Execute package type info map implementation
func (tim PkgsTiMap) Execute(
	ctx context.Context,
	st *types.Struct,
	fset *token.FileSet,
	tie TiExt,
) error {
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		ti := tie(field.Type())
		tim[field.Name()] = ti
	}
	return nil
}

// PkgsTiOut defines strategy implementation
// that uses pakage strategy type info map result
// and outputs it to io writer chanel
type PkgsTiOut struct {
	tim PkgsTiMap
	w   io.Writer
	f   func(interface{}) ([]byte, error)
}

// Execute package type info out implementation
func (tio PkgsTiOut) Execute(
	ctx context.Context,
	st *types.Struct,
	fset *token.FileSet,
	tie TiExt,
) error {
	if tio.f == nil {
		return errors.New("strategy type info out formatter method wasn't defined")
	}
	err := tio.tim.Execute(ctx, st, fset, tie)
	if err != nil {
		return err
	}
	buf, err := tio.f(tio.tim)
	if err != nil {
		return err
	}
	_, err = tio.w.Write(buf)
	return err
}
