package main

import (
	"context"
	"regexp"

	"1pkg/gopium/pkgs"
	"1pkg/gopium/pkgs/read"
	"1pkg/gopium/types"
	"1pkg/gopium/types/strategy"

	"golang.org/x/tools/go/packages"
)

// small gopium self example
func main() {
	// compile regex
	regex, _ := regexp.Compile(`.*`)
	// set up StrategyBuilder
	e := types.NewExtractorGoTypes("gc", "amd64")
	bs := strategy.NewBuilder(e)
	// build Strategy
	stg, err := bs.Build(strategy.StrategyEnumerate)
	if err != nil {
		panic(err)
	}
	// set up WalkerBuilder
	p := pkgs.ParserXToolPackages{
		Pattern:  "1pkg/gopium/pkgs",
		LoadMode: packages.LoadAllSyntax,
	}
	bw := read.NewBuilder(p)
	// build Walker
	w, err := bw.Build(read.WalkerOutPrettyJsonStd)
	if err != nil {
		panic(err)
	}
	// run VisitTop for Strategy with regex
	err = w.VisitTop(context.Background(), regex, stg)
	if err != nil {
		panic(err)
	}
}
