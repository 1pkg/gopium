package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/mocks"
)

func TestCache(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		cache   cache
		curator gopium.Curator
		ctx     context.Context
		o       gopium.Struct
		r       gopium.Struct
		err     error
	}{
		"empty struct should be applied to empty struct": {
			cache:   cachel1,
			curator: mocks.Maven{SCache: []int64{32}},
			ctx:     context.Background(),
		},
		"non empty struct should be applied to cache aligned struct": {
			cache:   cachel2,
			curator: mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					gopium.PadField(8),
				},
			},
		},
		"non empty struct should be applied to cache aligned struct on canceled context": {
			cache:   cachel3,
			curator: mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:     cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					gopium.PadField(8),
				},
			},
			err: cctx.Err(),
		},
		"mixed struct should be applied to cache aligned struct": {
			cache:   cachel3,
			curator: mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
					gopium.PadField(13),
				},
			},
		},
		"mixed prealigned struct should be applied to same cache aligned struct": {
			cache:   cachel2,
			curator: mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 16,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 16,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			cache := tcase.cache.Curator(tcase.curator)
			r, err := cache.Apply(tcase.ctx, tcase.o)
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
