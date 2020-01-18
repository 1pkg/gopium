package main

import (
	"go/ast"
	"regexp"
)

type Apply func(ast.StructType)

type Walker interface {
	Visit(reg *regexp.Regexp, apply Apply)
}
