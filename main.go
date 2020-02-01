package main

import (
	"context"
	"regexp"

	"golang.org/x/tools/go/packages"
)

func main() {
	// small pkg `regexp` example
	b := Pkgsb(GetTi)
	sg, err := b.Build(TypeInfoJsonStdOut)
	if err != nil {
		panic(err)
	}
	r, _ := regexp.Compile(`.*`)
	p := PkgpDef{
		Patterns: []string{"1pkg/gopium"},
		LoadMode: packages.LoadAllSyntax,
	}.Parse
	w, err := NewPackageWalker(context.Background(), r, p)
	if err != nil {
		panic(err)
	}
	w.VisitTop(context.Background(), r, sg)
}
