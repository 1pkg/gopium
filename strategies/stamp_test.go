package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestStamp(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		stamp stamp
		ctx   context.Context
		o     gopium.Struct
		r     gopium.Struct
		err   error
	}{
		"empty struct should be applied to itself with expected doc": {
			stamp: stampdoc,
			ctx:   context.Background(),
			r: gopium.Struct{
				Doc: []string{"// struct has been auto curated - 🌺 gopium @1pkg"},
			},
		},
		"non empty struct should be applied to itself with expected doc": {
			stamp: stampdoc,
			ctx:   context.Background(),
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
				Doc:  []string{"// struct has been auto curated - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
		},
		"non empty struct should be applied to itself with expected comment on canceled context": {
			stamp: stampcom,
			ctx:   cctx,
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
				Comment: []string{"// struct has been auto curated - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			err: context.Canceled,
		},
		"complex struct should be applied to itself with expected doc": {
			stamp: stampdoc,
			ctx:   context.Background(),
			o: gopium.Struct{
				Name: "test",
				Doc:  []string{"test"},
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Doc:  []string{"test", "// struct has been auto curated - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
		},
		"complex struct should be applied to itself with expected comment": {
			ctx: context.Background(),
			o: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test"},
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"test", "// struct has been auto curated - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test"},
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.stamp.Apply(tcase.ctx, tcase.o)
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
