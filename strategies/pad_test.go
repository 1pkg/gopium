package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/tests/mocks"
)

func TestPad(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		pad pad
		c   gopium.Curator
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			pad: padsys,
			c:   mocks.Maven{SAlign: 16},
			ctx: context.Background(),
		},
		"non empty struct should be applied to expected aligned struct": {
			pad: padsys,
			c:   mocks.Maven{SAlign: 6},
			ctx: context.Background(),
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
					collections.PadField(4),
				},
			},
		},
		"non empty struct should be applied to expected aligned struct on canceled context": {
			pad: padtnat,
			c:   mocks.Maven{SAlign: 12},
			ctx: cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 5,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 5,
					},
					collections.PadField(2),
				},
			},
			err: context.Canceled,
		},
		"mixed struct should be applied to expected aligned struct on type natural pad": {
			pad: padtnat,
			c:   mocks.Maven{SAlign: 24},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test4",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test1",
						Size:  32,
						Align: 4,
					},
					{
						Name:  "test2",
						Size:  6,
						Align: 6,
					},
					{
						Name:  "test3",
						Size:  8,
						Align: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test4",
						Size:  3,
						Align: 1,
					},
					collections.PadField(1),
					{
						Name:  "test1",
						Size:  32,
						Align: 4,
					},
					{
						Name:  "test2",
						Size:  6,
						Align: 6,
					},
					collections.PadField(6),
					{
						Name:  "test3",
						Size:  8,
						Align: 8,
					},
				},
			},
		},
		"mixed struct should be applied to expected aligned on field sys pad": {
			pad: padsys,
			c:   mocks.Maven{SAlign: 9},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 9,
					},
					{
						Name: "test2",
						Size: 7,
					},
					{
						Name: "test3",
						Size: 5,
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
						Size: 9,
					},
					{
						Name: "test2",
						Size: 7,
					},
					collections.PadField(2),
					{
						Name: "test3",
						Size: 5,
					},
					collections.PadField(4),
					{
						Name: "test4",
						Size: 3,
					},
					collections.PadField(6),
				},
			},
		},
		"mixed struct should be applied to expected aligned on big sys pad": {
			pad: padsys,
			c:   mocks.Maven{SAlign: 12},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 9,
					},
					{
						Name: "test2",
						Size: 7,
					},
					{
						Name: "test3",
						Size: 5,
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
						Size: 9,
					},
					collections.PadField(3),
					{
						Name: "test2",
						Size: 7,
					},
					collections.PadField(5),
					{
						Name: "test3",
						Size: 5,
					},
					collections.PadField(7),
					{
						Name: "test4",
						Size: 3,
					},
					collections.PadField(9),
				},
			},
		},
		"mixed struct should be applied to expected aligned no additional aligment": {
			pad: padsys,
			c:   mocks.Maven{SAlign: 4},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 24,
					},
					{
						Name: "test2",
						Size: 12,
					},
					{
						Name: "test3",
						Size: 36,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 24,
					},
					{
						Name: "test2",
						Size: 12,
					},
					{
						Name: "test3",
						Size: 36,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			pad := tcase.pad.Curator(tcase.c)
			// exec
			r, err := pad.Apply(tcase.ctx, tcase.o)
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
