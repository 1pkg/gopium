package astutil

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

	"golang.org/x/sync/errgroup"
)

func note(
	ctx context.Context,
	pkg *ast.Package,
	loc gopium.Locator,
	hsts HierarchyStructs,
	fsets map[string]*token.FileSet,
) (*ast.Package, error) {
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
	// concurently go through package files
	for name, file := range pkg.Files {
		name := name
		file := file
		ofile := file
		group.Go(func() error {
			// regenerate ast for file
			// in order to update all
			// definitions position
			// printing ast to buffer
			var buf bytes.Buffer
			err := printer.Fprint(
				&buf,
				loc.Fset(),
				ofile,
			)
			if err != nil {
				return err
			}
			// parser ast back to file
			fset := token.NewFileSet()
			if file, err = parser.ParseFile(
				fset,
				"",
				buf.String(),
				parser.ParseComments,
			); err != nil {
				return err
			}
			fsets[name] = fset
			// go through file structs
			// and note all commentss
			cat := loc.Cat(ofile.Pos())
			if file, err = walkFile(
				gctx,
				file,
				ordered(hsts[cat]),
				func(ts *ast.TypeSpec, st gopium.Struct) error {
					// check that we are working with ast.StructType
					tts, ok := ts.Type.(*ast.StructType)
					if !ok {
						return errors.New("notesync could only be applied to ast.StructType")
					}
					// prepare struct docs list
					stdocs := make([]*ast.Comment, 0, len(st.Doc))
					// collect all docs from resulted structure
					for _, doc := range st.Doc {
						// doc position is position of name - name len - 1
						slash := ts.Name.Pos() - token.Pos(len(ts.Name.Name)) - token.Pos(1)
						sdoc := ast.Comment{Slash: slash, Text: doc}
						stdocs = append(stdocs, &sdoc)
					}
					// update file docs list
					file.Comments = append(file.Comments, &ast.CommentGroup{List: stdocs})
					// prepare struct comments list
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
					// prepare fields storage for docs list and comments list
					fdocs := make([][]*ast.Comment, 0, len(st.Fields))
					fcoms := make([][]*ast.Comment, 0, len(st.Fields))
					// go through all resulted structure fields
					for _, field := range st.Fields {
						// collect all docs from resulted structure
						docs := make([]*ast.Comment, 0, len(field.Doc))
						for _, doc := range field.Doc {
							doc := ast.Comment{Text: doc}
							docs = append(docs, &doc)
						}
						// collect all comments from resulted structure
						coms := make([]*ast.Comment, 0, len(field.Comment))
						for _, com := range field.Comment {
							com := ast.Comment{Text: com}
							coms = append(coms, &com)
						}
						// put collected results to storages
						fdocs = append(fdocs, docs)
						fcoms = append(fcoms, coms)
					}
					// go through all original structure fields
					for index, field := range tts.Fields.List {
						// if we have docs in storage
						// append them to collected list
						if len(fdocs) > index {
							fdocs := fdocs[index]
							// set original slash pos
							for _, fdoc := range fdocs {
								// doc position is position of name - 1
								fdoc.Slash = field.Pos() - token.Pos(1)
							}
							// update file docs list
							file.Comments = append(file.Comments, &ast.CommentGroup{List: fdocs})
						}
						// if we have comments in storage
						// append them to collected list
						if len(fcoms) > index {
							fcoms := fcoms[index]
							// set original slash pos
							for _, fcom := range fcoms {
								// comment position is end of field type
								fcom.Slash = field.Type.End()
							}
							// update file comments list
							file.Comments = append(file.Comments, &ast.CommentGroup{List: fcoms})
						}
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
