package main

import (
	"go/ast"
	"regexp"
)

// Apply is custom callback type that applies some action on ast.StructType
type Apply func(*ast.StructType)

// Walker is interface that describes hierarchical walker that
// applies some action on ast.StructType
type Walker interface {
	Visit(reg *regexp.Regexp, apply Apply)
}
