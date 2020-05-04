package mocks

import (
	"bytes"
	"context"
	"go/ast"
	"go/build"
	"strings"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
)

// Persist defines mock astutil persist implementation
type Persist struct {
	Buffers map[string]bytes.Buffer
	Err     error
}

// Persist mock implementation
func (pr *Persist) Persist(ctx context.Context, p astutil.Print, pkg *ast.Package, loc gopium.Locator) error {
	// in case we have error
	// return it back
	if pr.Err != nil {
		return pr.Err
	}
	// go through each file one by one
	pr.Buffers = make(map[string]bytes.Buffer)
	for name, file := range pkg.Files {
		// prepare new buf for each file
		var buf bytes.Buffer
		// grab the latest file fset
		fset, _ := loc.Fset(name, nil)
		// write updated ast file to related os file
		// use original file set to keep format
		// in case any error happened
		// just return error back
		if err := p(&buf, fset, file); err != nil {
			return err
		}
		// patch the file name for tests
		name = strings.Replace(name, build.Default.GOPATH, "", 1)
		pr.Buffers[name] = buf
	}
	return nil
}
