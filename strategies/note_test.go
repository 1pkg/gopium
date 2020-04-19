package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestNoteDoc(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to itself with relevant doc": {
			ctx: context.Background(),
			r: gopium.Struct{
				Doc: []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
			},
		},
		"empty struct should be applied to itself with relevant doc on canceled context": {
			ctx: cctx,
			r: gopium.Struct{
				Doc: []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
			},
			err: cctx.Err(),
		},
		"non empty struct should be applied to itself with relevant doc": {
			ctx: context.Background(),
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
				Doc:  []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test",
						Doc:  []string{"// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct should be applied to itself with relevant doc on canceled context": {
			ctx: cctx,
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
				Doc:  []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test",
						Doc:  []string{"// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
			err: cctx.Err(),
		},
		"complex struct should be applied to itself with relevant doc": {
			ctx: context.Background(),
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
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Doc:  []string{"test", "// struct size: 16 bytes struct align: 8 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
						Doc:   []string{"// field size: 8 bytes field align: 4 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test", "// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Doc:   []string{"// field size: 8 bytes field align: 8 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"complex struct should be applied to itself with relevant doc on canceled context": {
			ctx: cctx,
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
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Doc:  []string{"test", "// struct size: 16 bytes struct align: 8 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
						Doc:   []string{"// field size: 8 bytes field align: 4 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test", "// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
						Doc:   []string{"// field size: 8 bytes field align: 8 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
			err: cctx.Err(),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := notedoc.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}

func TestNoteCom(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to itself with relevant doc": {
			ctx: context.Background(),
			r: gopium.Struct{
				Comment: []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
			},
		},
		"empty struct should be applied to itself with relevant doc on canceled context": {
			ctx: cctx,
			r: gopium.Struct{
				Comment: []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
			},
			err: cctx.Err(),
		},
		"non empty struct should be applied to itself with relevant doc": {
			ctx: context.Background(),
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
				Comment: []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:    "test",
						Comment: []string{"// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct should be applied to itself with relevant doc on canceled context": {
			ctx: cctx,
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
				Comment: []string{"// struct size: 0 bytes struct align: 0 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:    "test",
						Comment: []string{"// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
			err: cctx.Err(),
		},
		"complex struct should be applied to itself with relevant doc": {
			ctx: context.Background(),
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
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test", "// struct size: 16 bytes struct align: 8 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:    "test1",
						Type:    "int",
						Size:    8,
						Align:   4,
						Comment: []string{"// field size: 8 bytes field align: 4 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test", "// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "float64",
						Size:    8,
						Align:   8,
						Comment: []string{"// field size: 8 bytes field align: 8 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"complex struct should be applied to itself with relevant doc on canceled context": {
			ctx: cctx,
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
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test", "// struct size: 16 bytes struct align: 8 bytes - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:    "test1",
						Type:    "int",
						Size:    8,
						Align:   4,
						Comment: []string{"// field size: 8 bytes field align: 4 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test", "// field size: 0 bytes field align: 0 bytes - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test2",
						Type:    "float64",
						Size:    8,
						Align:   8,
						Comment: []string{"// field size: 8 bytes field align: 8 bytes - 🌺 gopium @1pkg"},
					},
				},
			},
			err: cctx.Err(),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := notecom.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
