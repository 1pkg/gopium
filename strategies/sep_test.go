package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/mocks"
)

func TestSep(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		sep     sep
		curator gopium.Curator
		ctx     context.Context
		o       gopium.Struct
		r       gopium.Struct
		err     error
	}{
		"empty struct should be applied to empty struct": {
			sep:     sepl1b,
			curator: mocks.Maven{SysCacheVals: []int64{32}},
			ctx:     context.Background(),
		},
		"non empty struct should be applied to cache line separator aligned struct": {
			sep:     sepl2b,
			curator: mocks.Maven{SysCacheVals: []int64{16, 16, 16}},
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
					gopium.PadField(16),
				},
			},
		},
		"non empty struct should be applied to cache line separator aligned on canceled context": {
			sep:     sepl3b,
			curator: mocks.Maven{SysCacheVals: []int64{16, 16, 16}},
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
					gopium.PadField(16),
				},
			},
			err: cctx.Err(),
		},
		"mixed struct should be applied to sys separator aligned struct at top": {
			sep:     sepsyst,
			curator: mocks.Maven{SysAlignVal: 24, SysCacheVals: []int64{16, 32, 64}},
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
					gopium.PadField(24),
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
		},
		"mixed struct should be applied to sys separator aligned struct at bottom": {
			sep:     sepsysb,
			curator: mocks.Maven{SysAlignVal: 24, SysCacheVals: []int64{16, 32, 64}},
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
					gopium.PadField(24),
				},
			},
		},
		"mixed struct should be applied to cache line separator aligned struct at top": {
			sep:     sepl3t,
			curator: mocks.Maven{SysAlignVal: 24, SysCacheVals: []int64{16, 32, 64}},
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
					{
						Name: "test5",
						Size: 1,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					gopium.PadField(64),
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
					{
						Name: "test5",
						Size: 1,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			sep := tcase.sep.Curator(tcase.curator)
			r, err := sep.Apply(tcase.ctx, tcase.o)
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
