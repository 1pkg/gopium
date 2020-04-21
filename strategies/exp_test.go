package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestExp(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		exp exp
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			exp: expasc,
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself": {
			exp: expasc,
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
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
		},
		"non empty struct should be applied to itself on canceled context": {
			exp: expdesc,
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
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			err: cctx.Err(),
		},
		"not exported struct should be applied to itself": {
			exp: expasc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
					{
						Name: "test3",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
					{
						Name: "test3",
					},
				},
			},
		},
		"exported struct should be applied to itself": {
			exp: expdesc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Exported: true,
					},
					{
						Name:     "test2",
						Exported: true,
					},
					{
						Name:     "test3",
						Exported: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Exported: true,
					},
					{
						Name:     "test2",
						Exported: true,
					},
					{
						Name:     "test3",
						Exported: true,
					},
				},
			},
		},
		"mixed exported struct should be applied to sorted struct asc": {
			exp: expasc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Exported: true,
					},
					{
						Name: "test2",
					},
					{
						Name:     "test3",
						Exported: true,
					},
					{
						Name: "test4",
					},
					{
						Name:     "test5",
						Exported: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Exported: true,
					},
					{
						Name:     "test3",
						Exported: true,
					},
					{
						Name:     "test5",
						Exported: true,
					},
					{
						Name: "test2",
					},
					{
						Name: "test4",
					},
				},
			},
		},
		"mixed exported struct should be applied to sorted struct desc": {
			exp: expdesc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Exported: true,
					},
					{
						Name: "test2",
					},
					{
						Name:     "test3",
						Exported: true,
					},
					{
						Name: "test4",
					},
					{
						Name:     "test5",
						Exported: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test2",
					},
					{
						Name: "test4",
					},
					{
						Name:     "test1",
						Exported: true,
					},
					{
						Name:     "test3",
						Exported: true,
					},
					{
						Name:     "test5",
						Exported: true,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.exp.Apply(tcase.ctx, tcase.o)
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
