package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/tests/mocks"
)

func TestCache(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		cache cache
		c     gopium.Curator
		ctx   context.Context
		o     gopium.Struct
		r     gopium.Struct
		err   error
	}{
		"empty struct should be applied to empty struct": {
			cache: cachel1,
			c:     mocks.Maven{SCache: []int64{32}},
			ctx:   context.Background(),
		},
		"non empty struct should be applied to expected aligned struct": {
			cache: cachel2,
			c:     mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:   context.Background(),
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
				},
			},
		},
		"non empty struct should be applied to expected aligned struct with full cahce": {
			cache: fcachel2,
			c:     mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:   context.Background(),
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
					collections.PadField(8),
				},
			},
		},
		"non empty struct should be applied to expected aligned struct on canceled context": {
			cache: fcachel3,
			c:     mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:   cctx,
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
					collections.PadField(8),
				},
			},
			err: context.Canceled,
		},
		"mixed struct should be applied to expected aligned struct": {
			cache: cachel3,
			c:     mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:   context.Background(),
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
					collections.PadField(13),
				},
			},
		},
		"mixed prealigned struct should be applied to itself": {
			cache: cachel2,
			c:     mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:   context.Background(),
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
		"struct with pads should be applied to expected aligned struct": {
			cache: cachel2,
			c:     mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:   context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Size:  8,
						Align: 8,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Size:  8,
						Align: 8,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
					collections.PadField(13),
				},
			},
		},
		"struct with explicit pads should be applied to expected aligned struct": {
			cache: cachel2,
			c:     mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:   context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					collections.PadField(5),
					{
						Name:  "test2",
						Size:  8,
						Align: 8,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
					collections.PadField(5),
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					collections.PadField(5),
					{
						Name:  "test2",
						Size:  8,
						Align: 8,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
					collections.PadField(5),
					collections.PadField(8),
				},
			},
		},
		"struct with pads should be applied to expected aligned struct div cache line": {
			cache: cachel2,
			c:     mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:   context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Size:  8,
						Align: 6,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Size:  8,
						Align: 6,
					},
					collections.PadField(2),
				},
			},
		},
		"struct with pads should be applied to expected aligned struct full cache line": {
			cache: fcachel2,
			c:     mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:   context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Size:  8,
						Align: 6,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Size:  8,
						Align: 6,
					},
					collections.PadField(18),
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			cache := tcase.cache.Curator(tcase.c)
			// exec
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
