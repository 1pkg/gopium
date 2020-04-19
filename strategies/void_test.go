package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestVoid(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			ctx: context.Background(),
		},
		"empty struct should be applied to empty struct on canceled context": {
			ctx: cctx,
			err: cctx.Err(),
		},
		"non empty struct should be applied to empty struct": {
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
		},
		"non empty struct should be applied to empty struct on canceled context": {
			ctx: cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			err: cctx.Err(),
		},
		"complex struct should be applied to empty struct": {
			ctx: context.Background(),
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
		},
		"complex struct should be applied to empty struct on canceled context": {
			ctx: cctx,
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
			err: cctx.Err(),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := vd.Apply(tcase.ctx, tcase.o)
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
