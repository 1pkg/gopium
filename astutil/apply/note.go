package apply

import (
	"bytes"
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"

	"1pkg/gopium"
	"1pkg/gopium/collections"

	"golang.org/x/sync/errgroup"
)

// note helps to update ast.Package
// accordingly to gopium.Struct result,
// it synchronizes all docs and comments by
//  regenerating ast for each file
// in order to update all definitions position
// and ingest docs and comments directly
// to file with correct calculated position
func note(
	ctx context.Context,
	pkg *ast.Package,
	loc gopium.Locator,
	hsts collections.Hierarchic,
) (*ast.Package, error) {
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
	// concurently go through package files
	for name, file := range pkg.Files {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-gctx.Done():
			return pkg, gctx.Err()
		default:
		}
		// capture name and file copies
		name := name
		file := file
		group.Go(func() error {
			// manage context actions
			// in case of cancelation
			// stop execution and return error
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}
			// print ast to buffer
			var buf bytes.Buffer
			err := printer.Fprint(
				&buf,
				loc.Root(),
				file,
			)
			if err != nil {
				return err
			}
			// parse ast back to file
			fset := token.NewFileSet()
			if file, err = parser.ParseFile(
				fset,
				"",
				buf.String(),
				parser.ParseComments,
			); err != nil {
				return err
			}
			// push child fset to locator
			loc.Fset(name, fset)
			// go through file structs
			// and note all comments
			if file, err = walkFile(
				gctx,
				file,
				compwnote(comploc(loc, name, hsts)),
				func(ts *ast.TypeSpec, st gopium.Struct) error {
					// check that we are working with ast.StructType
					tts, ok := ts.Type.(*ast.StructType)
					if !ok {
						return errors.New("notesync could only be applied to ast.StructType")
					}
					// prepare struct docs slice
					stdocs := make([]*ast.Comment, 0, len(st.Doc))
					// collect all docs from resulted structure
					for _, doc := range st.Doc {
						// doc position is position of name - name len - 1
						slash := ts.Name.Pos() - token.Pos(len(ts.Name.Name)) - token.Pos(1)
						sdoc := ast.Comment{Slash: slash, Text: doc}
						stdocs = append(stdocs, &sdoc)
					}
					// update file comments list
					file.Comments = append(file.Comments, &ast.CommentGroup{List: stdocs})
					// prepare struct comments slice
					stcoms := make([]*ast.Comment, 0, len(st.Comment))
					// collect all comments from resulted structure
					for _, com := range st.Comment {
						// comment position is end of type decl
						slash := ts.Type.End()
						scom := ast.Comment{Slash: slash, Text: com}
						stcoms = append(stcoms, &scom)
					}
					// update file comments list
					file.Comments = append(file.Comments, &ast.CommentGroup{List: stcoms})
					// go through all resulted structure fields
					for index, field := range st.Fields {
						// if index is greater that ast
						// field number break the loop
						if len(tts.Fields.List) <= index {
							break
						}
						// get the field from ast
						astfield := tts.Fields.List[index]
						// collect all docs from resulted structure
						fdocs := make([]*ast.Comment, 0, len(field.Doc))
						for _, doc := range field.Doc {
							// doc position is position of name - 1
							slash := astfield.Pos() - token.Pos(1)
							fdoc := ast.Comment{Slash: slash, Text: doc}
							fdocs = append(fdocs, &fdoc)
						}
						// update file comments list
						file.Comments = append(file.Comments, &ast.CommentGroup{List: fdocs})
						// collect all comments from resulted structure
						fcoms := make([]*ast.Comment, 0, len(field.Comment))
						for _, com := range field.Comment {
							// comment position is end of field type
							slash := astfield.Type.End()
							fcom := ast.Comment{Slash: slash, Text: com}
							fcoms = append(fcoms, &fcom)
						}
						// update file comments list
						file.Comments = append(file.Comments, &ast.CommentGroup{List: fcoms})
					}
					return nil
				},
			); err != nil {
				return err
			}
			// sort all comments by their ast pos
			sort.SliceStable(file.Comments, func(i, j int) bool {
				// sort safe guard
				if len(file.Comments[i].List) == 0 {
					return true
				}
				// sort safe guard
				if len(file.Comments[j].List) == 0 {
					return false
				}
				return file.Comments[i].Pos() < file.Comments[j].Pos()
			})
			// update result package file
			pkg.Files[name] = file
			return nil
		})
	}
	// wait until walk is done
	return pkg, group.Wait()
}
