package main

import (
	"context"
	"go/parser"
	"path/filepath"
	"regexp"

	"1pkg/gopium/pkgs_types"
	"1pkg/gopium/strategy"
	"1pkg/gopium/walker"

	"golang.org/x/tools/go/packages"
)

// small gopium self example
func main() {
	// set up StrategyBuilder
	m := pkgs_types.NewMavenGoTypes("gc", "amd64")
	bs := strategy.NewBuilder(m)
	// build Strategy
	stg, err := bs.Build(strategy.LexAsc)
	if err != nil {
		panic(err)
	}
	// set up WalkerBuilder
	abs, err := filepath.Abs("./../../../..")
	if err != nil {
		panic(err)
	}
	p := pkgs_types.ParserXToolPackagesAst{
		Pattern: "1pkg/gopium/pkgs_types",
		AbsDir:  abs,
		//nolint
		ModeTypes: packages.LoadAllSyntax,
		ModeAst:   parser.ParseComments,
	}
	bw := walker.NewBuilder(p, m, false)
	// build Walker
	w, err := bw.Build(walker.PrettyJsonStd)
	if err != nil {
		panic(err)
	}
	// run VisitTop for Strategy with regex
	err = w.VisitDeep(
		context.Background(),
		regexp.MustCompile(`.*`),
		stg,
	)
	if err != nil {
		panic(err)
	}
}
