package astutil

import (
	"bytes"
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"sync"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/fmtio"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/typepkg"

	"golang.org/x/sync/errgroup"
)

// UFFN implements apply and combines:
// - ufmt with fmtio FSPT helper
// - filter helper
// - note helper
var UFFN = combine(
	ufmt(walk, fmtio.FSPT),
	filter(walk),
	note(
		walk,
		&typepkg.ParserXToolPackagesAst{
			ModeAst: parser.ParseComments | parser.AllErrors,
		},
		fmtio.Gofmt{},
	),
)

// combine helps to pipe several
// ast helpers to single apply func
func combine(funcs ...gopium.Apply) gopium.Apply {
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		c gopium.Categorized,
	) (rpkg *ast.Package, err error) {
		// filters all files that
		// could be skipped
		// we could skip err check here
		// as cat alway returns nil
		pkg, _ = cat(ctx, pkg, loc, c)
		// go through all provided funcs
		for _, fun := range funcs {
			// manage context actions
			// in case of cancelation
			// stop execution
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			// exec single func
			pkg, err = fun(ctx, pkg, loc, c)
			// in case of any error
			// just propagate it
			if err != nil {
				return nil, err
			}
		}
		return pkg, nil
	}
}

// cat helps to filter only ast files
// that exist inside hierarchic collection
func cat(
	_ context.Context,
	//nolint
	pkg *ast.Package,
	_ gopium.Locator,
	c gopium.Categorized,
) (*ast.Package, error) {
	files := make(map[string]*ast.File, len(pkg.Files))
	for name, file := range pkg.Files {
		// skip empty writes
		if _, ok := c.Cat(name); ok {
			files[name] = file
		}
	}
	pkg.Files = files
	return pkg, nil
}

// ufmt helps to update ast package
// accordingly to gopium struct result
// using custom fmtio ast formatter
func ufmt(w gopium.Walk, fmt gopium.Ast) gopium.Apply {
	//nolint
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		c gopium.Categorized,
	) (*ast.Package, error) {
		// walk through the ast
		// and use fmt to nodes
		wpkg, err := w(
			ctx,
			pkg,
			fmtast(fmt),
			&flatid{loc: loc, sts: collections.Flat(c.Full())},
		)
		if err != nil {
			return nil, err
		}
		// if no error happened
		// just apply format to ast
		return wpkg.(*ast.Package), nil
	}
}

// filter helps to filter structs
// docs and comments from ast type spec
//
// it filters only comments inside
// result structs and autogenerated comments
func filter(w gopium.Walk) gopium.Apply {
	//nolint
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		c gopium.Categorized,
	) (*ast.Package, error) {
		// prepare structs boundaries
		bc := &bcollect{}
		// collect structs boundaries
		if _, err := w(
			ctx,
			pkg,
			bc,
			&flatid{loc: loc, sts: collections.Flat(c.Full())},
		); err != nil {
			return nil, err
		}
		// create sync error group
		// with cancelation context
		group, gctx := errgroup.WithContext(ctx)
		// go through package files
		for _, file := range pkg.Files {
			// manage context actions
			// in case of cancelation
			// stop execution
			select {
			case <-gctx.Done():
				break
			default:
			}
			// capture file copy
			file := file
			group.Go(func() error {
				// go through all file comments
				for _, comments := range file.Comments {
					// prepare comment slice
					comlist := make([]*ast.Comment, 0, len(comments.List))
					// go through comment slice
					for _, com := range comments.List {
						// if comment is inside boundaries skip it
						if bc.bs.Inside(com.Slash) {
							continue
						}
						// if comment has autogenerated stamp
						// and it locates between struct's
						// start and end points skip it
						if strings.Contains(com.Text, gopium.STAMP) &&
							(bc.bs.Inside(com.Slash-1) || bc.bs.Inside(com.Slash+token.Pos(len(com.Text)+1))) {
							continue
						}
						// otherwise append comment to slice
						comlist = append(comlist, com)
					}
					// update comment list
					comments.List = comlist
				}
				return gctx.Err()
			})
		}
		// wait until walk is done
		// in case of any error
		// just return it back
		if err := group.Wait(); err != nil {
			return nil, err
		}
		return pkg, nil
	}
}

// note helps to update ast package
// accordingly to gopium struct result
//
// it synchronizes all docs and comments by
// regenerating ast for each file in order
// to update all definitions position
// and ingest docs and comments directly
// to file with correct calculated positions
func note(w gopium.Walk, xp gopium.AstParser, p gopium.Printer) gopium.Apply {
	//nolint
	return func(
		ctx context.Context,
		pkg *ast.Package,
		loc gopium.Locator,
		c gopium.Categorized,
	) (*ast.Package, error) {
		// create sync error group
		// with cancelation context
		// and sync map to store updated files
		var files sync.Map
		group, gctx := errgroup.WithContext(ctx)
		// concurently go through package files
		for name, file := range pkg.Files {
			// manage context actions
			// in case of cancelation
			// stop execution
			select {
			case <-gctx.Done():
				break
			default:
			}
			// capture name and file copies
			name := name
			file := file
			group.Go(func() error {
				// print ast to buffer
				var buf bytes.Buffer
				if err := p.Print(gctx, &buf, loc.Root(), file); err != nil {
					return err
				}
				// parse ast back to file
				// and push child fset to locator
				pkg, nloc, err := xp.ParseAst(gctx, buf.Bytes()...)
				if err != nil {
					return err
				}
				// get collection for cat
				// empty cat should be skipped
				// before in cat
				catsts, _ := c.Cat(name)
				// go through file structs
				// and note all comments
				file := pkg.Files["file"]
				node, err := w(
					gctx,
					file,
					((*pressnote)(file)),
					hasnote{cmp: newsorted(catsts)},
				)
				if err != nil {
					return err
				}
				// sort all comments by their ast pos
				file = node.(*ast.File)
				sort.SliceStable(file.Comments, func(i, j int) bool {
					return file.Comments[i].Pos() < file.Comments[j].Pos()
				})
				// save update file and loc results
				files.Store(name, file)
				loc.Fset(name, nloc.Root())
				return gctx.Err()
			})
		}
		// wait until walk is done
		// in case of any error
		// just return it back
		if err := group.Wait(); err != nil {
			return nil, err
		}
		// otherwise update ast pkg files
		// with synced files result
		files.Range(func(name interface{}, file interface{}) bool {
			pkg.Files[name.(string)] = file.(*ast.File)
			return true
		})
		return pkg, nil
	}
}
