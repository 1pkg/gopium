package main

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPackageWalker(t *testing.T) {
	table := []struct {
		name string
		pkgp Pkgp
		pkgw *Pkgw
		err  error
	}{
		{
			name: "package parser returns error, new package walker should just pass it",
			pkgp: PkgpErr("error package test error").Parse,
			pkgw: nil,
			err:  errors.New("error package test error"),
		},
		{
			name: "package name wasn't found, new package walker should return error",
			pkgp: PkgpNF{}.Parse,
			pkgw: nil,
			err:  errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name: "package parser returns nil type package, new package walker should return error",
			pkgp: PkgpMock{fset: token.NewFileSet()}.Parse,
			pkgw: nil,
			err:  errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name: "package parser returns nil fset, new package walker should return error",
			pkgp: PkgpMock{
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
			}.Parse,
			pkgw: nil,
			err:  errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name: "package was found, new package walker should return correct package walker",
			pkgp: PkgpMock{
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
				fset: token.NewFileSet(),
			}.Parse,
			pkgw: &Pkgw{
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
				fset: token.NewFileSet(),
			},
			err: nil,
		},
	}
	t.Run("new package walker should return correct results for all cases", func(t *testing.T) {
		r, _ := regexp.Compile(`^foobar$`)
		for _, tcase := range table {
			t.Run(tcase.name, func(t *testing.T) {
				ctx := context.Background()
				pkgw, err := NewPackageWalker(ctx, r, tcase.pkgp)
				assert.Equal(t, tcase.pkgw, pkgw)
				assert.Equal(t, tcase.err, err)
			})
		}
	})
}

// TODO add pkgw visit tests
