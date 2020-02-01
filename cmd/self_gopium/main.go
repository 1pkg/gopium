package main

import (
	"context"
	"regexp"

	"1pkg/gopium/pkgs"
	"1pkg/gopium/typeinfo"
	"1pkg/gopium/typeinfo/strategy"

	"golang.org/x/tools/go/packages"
)

// small gopium self example
func main() {
	// compile regex
	regex, _ := regexp.Compile(`.*`)
	// set up strategy builder
	e := typeinfo.NewExtractorTypesSizes("gc", "amd64")
	b := strategy.NewBuilder(e)
	stg, err := b.Build(strategy.TypeInfoOutPrettyJsonStd)
	if err != nil {
		panic(err)
	}
	// create pkgs Walker
	p := pkgs.ParserXTool{
		Patterns: []string{"1pkg/gopium/pkgs"},
		LoadMode: packages.LoadAllSyntax,
	}.Parse
	w, err := pkgs.NewWalker(context.Background(), regex, p)
	if err != nil {
		panic(err)
	}
	// run VisitTop for Strategy with regex
	w.VisitTop(context.Background(), regex, stg)
}
