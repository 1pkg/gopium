package main

import (
	"errors"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPackageWalker(t *testing.T) {
	table := []struct {
		name string
		fpkg string
		pkgp Pkgp
		pkgw *Pkgw
		err  error
	}{
		{
			name: "package parser returns error, new package walker should just pass it",
			fpkg: "foobar",
			pkgp: ErrorPkgp("error package test error").Parse,
			pkgw: nil,
			err:  errors.New("error package test error"),
		},
		{
			name: "empty package name, new package walker should return error",
			fpkg: "",
			pkgp: MockPkgp(map[string]*ast.Package{"foo": nil, "bar": nil}).Parse,
			pkgw: nil,
			err:  errors.New("package `` wasn't found"),
		},
		{
			name: "package parser returns nil map, new package walker should return error",
			fpkg: "foobar",
			pkgp: MockPkgp(nil).Parse,
			pkgw: nil,
			err:  errors.New("package `foobar` wasn't found"),
		},
		{
			name: "package parser returns empty map, new package walker should return error",
			fpkg: "foobar",
			pkgp: MockPkgp(make(map[string]*ast.Package)).Parse,
			pkgw: nil,
			err:  errors.New("package `foobar` wasn't found"),
		},
		{
			name: "package name wasn't found, new package walker should return error",
			fpkg: "foobar",
			pkgp: MockPkgp(map[string]*ast.Package{"foo": nil, "bar": nil}).Parse,
			pkgw: nil,
			err:  errors.New("package `foobar` wasn't found"),
		},
		{
			name: "single nil package was found, new package walker should return correct package walker",
			fpkg: "foobar",
			pkgp: MockPkgp(map[string]*ast.Package{"foobar": nil}).Parse,
			pkgw: &Pkgw{},
			err:  nil,
		},
		{
			name: "package was found in the map, new package walker should return correct package walker",
			fpkg: "foobar",
			pkgp: MockPkgp(map[string]*ast.Package{
				"foo":    &ast.Package{Name: "foo"},
				"bar":    &ast.Package{Name: "bar"},
				"foobar": &ast.Package{Name: "foobar"},
			}).Parse,
			pkgw: &Pkgw{pkg: &ast.Package{Name: "foobar"}},
			err:  nil,
		},
	}
	t.Run("new package walker should return correct results for all cases", func(t *testing.T) {
		for _, tcase := range table {
			t.Run(tcase.name, func(t *testing.T) {
				pkgw, err := NewPackageWalker(tcase.fpkg, tcase.pkgp)
				assert.Equal(t, tcase.pkgw, pkgw)
				assert.Equal(t, tcase.err, err)
			})
		}
	})
}

// TODO add pkgw visit tests
