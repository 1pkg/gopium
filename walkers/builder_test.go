package walkers

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests/mocks"
)

func TestBuilder(t *testing.T) {
	// prepare
	b := Builder{
		Parser:  mocks.Parser{},
		Exposer: mocks.Maven{},
		Deep:    true,
		Bref:    true,
	}
	table := map[string]struct {
		name gopium.WalkerName
		w    gopium.Walker
		err  error
	}{
		// wast walkers
		"`ast_std` name should return expected walker": {
			name: AstStd,
			w: aststd.With(
				b.Parser,
				b.Exposer,
				b.Printer,
				b.Deep,
				b.Bref,
			),
		},
		"`ast_go` name should return expected walker": {
			name: AstGo,
			w: astgo.With(
				b.Parser,
				b.Exposer,
				b.Printer,
				b.Deep,
				b.Bref,
			),
		},
		"`ast_go_tree` name should return expected walker": {
			name: AstGoTree,
			w: astgotree.With(
				b.Parser,
				b.Exposer,
				b.Printer,
				b.Deep,
				b.Bref,
			),
		},
		"`ast_gopium` name should return expected walker": {
			name: AstGopium,
			w: astgopium.With(
				b.Parser,
				b.Exposer,
				b.Printer,
				b.Deep,
				b.Bref,
			),
		},
		// wout walkers
		"`file_json` name should return expected walker": {
			name: FileJsonb,
			w: filejson.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`file_xml` name should return expected walker": {
			name: FileXmlb,
			w: filexml.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`file_csv` name should return expected walker": {
			name: FileCsvb,
			w: filecsv.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`file_md_table` name should return expected walker": {
			name: FileMdt,
			w: filemdt.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		// wdiff walkers
		"`size_align_file_md_table` name should return expected walker": {
			name: SizeAlignFileMdt,
			w: safilemdt.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`fields_file_html_table` name should return expected walker": {
			name: FieldsFileHtmlt,
			w: ffilehtml.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		// others
		"invalid name should return builder error": {
			name: "test",
			err:  fmt.Errorf(`walker "test" wasn't found`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			w, err := b.Build(tcase.name)
			// check
			// we can't compare functions directly in go
			// so apply this hack to compare with nil
			if tcase.w != nil && reflect.DeepEqual(w, nil) {
				t.Errorf("actual %v doesn't equal to expected %v", w, tcase.w)
			}
			if tcase.w == nil && !reflect.DeepEqual(w, nil) {
				t.Errorf("actual %v doesn't equal to expected not %v", w, tcase.w)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
