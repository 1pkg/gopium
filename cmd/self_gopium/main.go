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
	// compile regex
	regex, _ := regexp.Compile(`.*`)
	// set up StrategyBuilder
	wb := pkgs_types.NewMavenGoTypes("gc", "amd64")
	bs := strategy.NewBuilder(wb)
	// build Strategy
	stg, err := bs.Build(strategy.Lexicographical, strategy.WithAnnotation)
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
	bw := walker.NewBuilder(p)
	// build Walker
	w, err := bw.Build(walker.PrettyJsonStd)
	if err != nil {
		panic(err)
	}
	// run VisitTop for Strategy with regex
	err = w.VisitDeep(context.Background(), regex, stg)
	if err != nil {
		panic(err)
	}
}
