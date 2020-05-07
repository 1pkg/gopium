package walkers

import (
	"fmt"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/tests/mocks"
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
		"`json_std` name should return expected walker": {
			name: JsonStd,
			w: jsonstd.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`xml_std` name should return expected walker": {
			name: XmlStd,
			w: xmlstd.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`csv_std` name should return expected walker": {
			name: CsvStd,
			w: csvstd.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`json_files` name should return expected walker": {
			name: JsonFiles,
			w: jsonfiles.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`xml_files` name should return expected walker": {
			name: XmlFiles,
			w: xmlfiles.With(
				b.Parser,
				b.Exposer,
				b.Deep,
				b.Bref,
			),
		},
		"`csv_files` name should return expected walker": {
			name: CsvFiles,
			w: csvfiles.With(
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
