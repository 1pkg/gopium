package main

import (
	"context"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/pkgs"

	"golang.org/x/tools/go/packages"
)

func main() {
	// small pkg `regexp` example
	b := gopium.Pkgsb(gopium.GetTi)
	sg, err := b.Build(gopium.TypeInfoJsonStdOut)
	if err != nil {
		panic(err)
	}
	r, _ := regexp.Compile(`.*`)
	p := pkgs.ParserXTool{
		Patterns: []string{"1pkg/gopium"},
		LoadMode: packages.LoadAllSyntax,
	}.Parse
	w, err := pkgs.NewWalker(context.Background(), r, p)
	if err != nil {
		panic(err)
	}
	w.VisitTop(context.Background(), r, sg)
}
