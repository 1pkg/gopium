package pkgs

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
		p    Parser
		w    *Walker
		err  error
	}{
		{
			name: "package parser returns error, new package walker should just pass it",
			p:    ParserErr("error package test error").Parse,
			w:    nil,
			err:  errors.New("error package test error"),
		},
		{
			name: "package name wasn't found, new package walker should return error",
			p:    ParserNil{}.Parse,
			w:    nil,
			err:  errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name: "package parser returns nil type package, new package walker should return error",
			p:    ParserMock{fset: token.NewFileSet()}.Parse,
			w:    nil,
			err:  errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name: "package parser returns nil fset, new package walker should return error",
			p:    ParserMock{pkgs: []*types.Package{types.NewPackage("/", "foobar")}}.Parse,
			w:    nil,
			err:  errors.New(`packages "^foobar$" wasn't found`),
		},
		{
			name: "package was found, new package walker should return correct package walker",
			p: ParserMock{
				pkgs: []*types.Package{types.NewPackage("/", "foobar")},
				fset: token.NewFileSet(),
			}.Parse,
			w: &Walker{
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
				w, err := NewWalker(ctx, r, tcase.p)
				assert.Equal(t, tcase.w, w)
				assert.Equal(t, tcase.err, err)
			})
		}
	})
}

// TODO add walker visits tests
