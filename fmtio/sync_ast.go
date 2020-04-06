package fmtio

import (
	"fmt"
	"go/ast"
	"sort"

	"1pkg/gopium"

	"golang.org/x/tools/go/ast/astutil"
)

// SyncAst helps to update ast.Package
// accordingly to gopium.Struct result
// synchronously or return error otherwise
func SyncAst(
	pkg *ast.Package,
	loc gopium.Locator,
	id string,
	st gopium.Struct,
	sta StructToAst,
) (*ast.Package, error) {
	// tracks error inside astutil.Apply
	var err error
	// apply astutil.Apply to parsed ast.Package
	// and update structure in ast
	snode := astutil.Apply(pkg, func(c *astutil.Cursor) bool {
		if gendecl, ok := c.Node().(*ast.GenDecl); ok {
			for _, spec := range gendecl.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := ts.Type.(*ast.StructType); ok {
						// calculate sum for structure
						// and skip all irrelevant structs
						lid := loc.ID(ts.Pos())
						if lid == id {
							// apply format to ast
							err = sta(ts, st)
							// in case we have error
							// break iteration
							return err != nil
						}
					}
				}
			}
		}
		return true
	}, nil)
	// in case we had error in astutil.Apply
	// just return it back
	if err != nil {
		return nil, err
	}
	// check that updated type is correct
	if spkg, ok := snode.(*ast.Package); ok {
		// go through all files and press
		// comment and doc to then
		for _, file := range spkg.Files {
			if f := comdocfpress(file); f == nil {
				// in case updated type isn't expected
				return nil, fmt.Errorf("can't update ast for structure %q", st.Name)
			}
		}
		return spkg, nil
	}
	// in case updated type isn't expected
	return nil, fmt.Errorf("can't update ast for structure %q", st.Name)
}

// comdocfpress helps to press ast.TypeSpec
// docs and comments directly to ast.File
// as printer only supports those comments
func comdocfpress(file *ast.File) *ast.File {
	// apply astutil.Apply to parsed ast.File
	// and press comment to file instead of struct
	pnode := astutil.Apply(file, func(c *astutil.Cursor) bool {
		if gendecl, ok := c.Node().(*ast.GenDecl); ok {
			for _, spec := range gendecl.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						// megre doc to file comments list
						if ts.Doc != nil {
							file.Comments = append(file.Comments, ts.Doc)
							ts.Doc = nil
						}
						// go through fields list
						for _, f := range st.Fields.List {
							// megre docs to file comments list
							if f.Doc != nil {
								file.Comments = append(file.Comments, f.Doc)
								f.Doc = nil
							}
							// megre comments to file comments list
							if ts.Comment != nil {
								file.Comments = append(file.Comments, f.Comment)
								f.Comment = nil
							}
						}
						// megre comment to file comments list
						if ts.Comment != nil {
							file.Comments = append(file.Comments, ts.Comment)
							ts.Comment = nil
						}
					}
				}
			}
		}
		return true
	}, nil)
	// check that updated type is correct
	if pfile, ok := pnode.(*ast.File); ok {
		// sort all comments by their ast pos
		sort.SliceStable(pfile.Comments, func(i, j int) bool {
			// sort safe guard
			if len(pfile.Comments[i].List) == 0 {
				return true
			}
			// sort safe guard
			if len(pfile.Comments[j].List) == 0 {
				return false
			}
			return pfile.Comments[i].Pos() < pfile.Comments[j].Pos()
		})
		return pfile
	}
	// in case updated type isn't expected
	return nil
}
