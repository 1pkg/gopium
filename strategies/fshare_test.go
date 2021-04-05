package strategies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests/mocks"
)

func TestFshare(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		fshare fshare
		c      gopium.Curator
		ctx    context.Context
		o      gopium.Struct
		r      gopium.Struct
		err    error
	}{
		"empty struct should be applied to empty struct": {
			fshare: fsharel1,
			c:      mocks.Maven{SCache: []int64{32}},
			ctx:    context.Background(),
		},
		"non empty struct should be applied to expected aligned struct": {
			fshare: fsharel2,
			c:      mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:    context.Background(),
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
			fshare: fsharel3,
			c:      mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:    cctx,
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
		"non empty struct should be applied to expected aligned struct custom bytes": {
			fshare: fshareb.Bytes(20),
			c:      mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:    context.Background(),
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
					collections.PadField(12),
				},
			},
		},
		"non empty struct should be applied to expected aligned struct custom bytes and cache line": {
			fshare: fsharel2.Bytes(20),
			c:      mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:    context.Background(),
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
		"non empty struct should be applied to expected aligned struct empty custom bytes and empty line": {
			fshare: fshareb,
			c:      mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx:    context.Background(),
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
		"mixed struct should be applied to expected aligned struct": {
			fshare: fsharel3,
			c:      mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:    context.Background(),
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
					collections.PadField(32),
					{
						Name: "test2",
						Size: 8,
					},
					collections.PadField(56),
					{
						Name: "test3",
						Size: 8,
					},
					collections.PadField(56),
					{
						Name: "test4",
						Size: 3,
					},
					collections.PadField(61),
				},
			},
		},
		"mixed struct should be applied to expected aligned struct custom bytes": {
			fshare: fshareb.Bytes(5),
			c:      mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:    context.Background(),
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
					collections.PadField(3),
					{
						Name: "test2",
						Size: 8,
					},
					collections.PadField(2),
					{
						Name: "test3",
						Size: 8,
					},
					collections.PadField(2),
					{
						Name: "test4",
						Size: 3,
					},
					collections.PadField(2),
				},
			},
		},
		"mixed prealigned struct should be applied to expected aligned struct": {
			fshare: fsharel1,
			c:      mocks.Maven{SCache: []int64{16, 32, 64}},
			ctx:    context.Background(),
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
					collections.PadField(8),
					{
						Name: "test3",
						Size: 8,
					},
					collections.PadField(8),
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			fshare := tcase.fshare.Curator(tcase.c)
			// exec
			r, err := fshare.Apply(tcase.ctx, tcase.o)
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
