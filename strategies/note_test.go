package strategies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
)

func TestNote(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		note note
		ctx  context.Context
		o    gopium.Struct
		r    gopium.Struct
		err  error
	}{
		"empty struct should be applied to itself": {
			note: fnotedoc,
			ctx:  context.Background(),
			r:    gopium.Struct{},
		},
		"non empty struct should be applied to itself with expected doc fields": {
			note: fnotedoc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Doc:  []string{"// field size: 0 bytes; field align: 0 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct should be applied to itself with expected comment fields on canceled context": {
			note: fnotecom,
			ctx:  cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:    "test",
						Comment: []string{"// field size: 0 bytes; field align: 0 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
			err: context.Canceled,
		},
		"empty struct should be applied to itself with expected doc struct": {
			note: stnotedoc,
			ctx:  context.Background(),
			r: gopium.Struct{
				Doc: []string{"// struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg"},
			},
		},
		"non empty struct should be applied to itself with expected doc struct": {
			note: stnotedoc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Doc:  []string{"// struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
		},
		"non empty struct should be applied to itself with expected comment struct on canceled context": {
			note: stnotecom,
			ctx:  cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"// struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; struct ptr scan size: 0 bytes; - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			err: context.Canceled,
		},
		"complex struct should be applied to itself with expected doc fields": {
			note: fnotedoc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Doc:  []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
						Ptr:   4,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   6,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Doc:  []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
						Ptr:   4,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 4 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test", "// field size: 0 bytes; field align: 0 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   6,
						Doc:   []string{"// field size: 8 bytes; field align: 8 bytes; field ptr: 6 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"complex struct should be applied to itself with expected comment fields": {
			note: fnotecom,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
						Ptr:   4,
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   6,
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:    "test1",
						Type:    "int",
						Size:    8,
						Align:   4,
						Ptr:     4,
						Comment: []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 4 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test", "// field size: 0 bytes; field align: 0 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "float64",
						Size:    8,
						Align:   8,
						Ptr:     6,
						Comment: []string{"// field size: 8 bytes; field align: 8 bytes; field ptr: 6 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"complex struct with pads should be applied to itself with expected comment fields": {
			note: fnotecom,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   4,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:    "test1",
						Size:    3,
						Align:   1,
						Comment: []string{"// field size: 3 bytes; field align: 1 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "float64",
						Size:    8,
						Align:   8,
						Ptr:     4,
						Comment: []string{"// field size: 8 bytes; field align: 8 bytes; field ptr: 4 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test3",
						Size:    3,
						Align:   1,
						Comment: []string{"// field size: 3 bytes; field align: 1 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"complex struct should be applied to itself with expected doc struct": {
			note: stnotedoc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Doc:  []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   2,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Doc:  []string{"test", "// struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; struct ptr scan size: 10 bytes; - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   2,
					},
				},
			},
		},
		"complex struct should be applied to itself with expected comment struct": {
			note: stnotecom,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   2,
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test", "// struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; struct ptr scan size: 10 bytes; - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   2,
					},
				},
			},
		},
		"complex struct with pads should be applied to itself with expected comment struct": {
			note: stnotecom,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   4,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test", "// struct size: 14 bytes; struct align: 8 bytes; struct aligned size: 24 bytes; struct ptr scan size: 7 bytes; - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Ptr:   4,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.note.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if !errors.Is(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
