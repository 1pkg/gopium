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
	p := PkgpDef{
		LoadMode: packages.LoadAllSyntax,
	}.Parse
	w, err := NewPackageWalker(context.Background(), "regexp", p)
	if err != nil {
		panic(err)
	}
	r, err := regexp.Compile(`.*`)
	if err != nil {
		panic(err)
	}
	w.VisitTop(context.Background(), r, sg)
}
