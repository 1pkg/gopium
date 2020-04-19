package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestEmbAsc(t *testing.T) {
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
		"non empty struct should be applied to itself": {
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
		"not embedded struct should be applied to itself": {
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
		"mixed embedded struct should be applied to sorted struct": {
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
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := embasc.Apply(tcase.ctx, tcase.o)
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

func TestEmbDesc(t *testing.T) {
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
		"non empty struct should be applied to itself": {
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
		"not embedded struct should be applied to itself": {
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
		"mixed embedded struct should be applied to sorted struct": {
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
			r, err := embdesc.Apply(tcase.ctx, tcase.o)
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
