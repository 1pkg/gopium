package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestEmb(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		emb emb
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			emb: embasc,
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself": {
			emb: embasc,
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
			emb: embdesc,
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
			err: context.Canceled,
		},
		"not embedded struct should be applied to itself": {
			emb: embasc,
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
		"embedded struct should be applied to itself": {
			emb: embdesc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Embedded: true,
					},
					{
						Name:     "test2",
						Embedded: true,
					},
					{
						Name:     "test3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Embedded: true,
					},
					{
						Name:     "test2",
						Embedded: true,
					},
					{
						Name:     "test3",
						Embedded: true,
					},
				},
			},
		},
		"mixed embedded struct should be applied to sorted struct asc": {
			emb: embasc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Embedded: true,
					},
					{
						Name: "test2",
					},
					{
						Name:     "test3",
						Embedded: true,
					},
					{
						Name: "test4",
					},
					{
						Name:     "test5",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Embedded: true,
					},
					{
						Name:     "test3",
						Embedded: true,
					},
					{
						Name:     "test5",
						Embedded: true,
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
		"mixed embedded struct should be applied to sorted struct desc": {
			emb: embdesc,
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test1",
						Embedded: true,
					},
					{
						Name: "test2",
					},
					{
						Name:     "test3",
						Embedded: true,
					},
					{
						Name: "test4",
					},
					{
						Name:     "test5",
						Embedded: true,
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
						Embedded: true,
					},
					{
						Name:     "test3",
						Embedded: true,
					},
					{
						Name:     "test5",
						Embedded: true,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.emb.Apply(tcase.ctx, tcase.o)
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
