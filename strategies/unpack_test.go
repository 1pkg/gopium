package strategies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
)

func TestUnpack(t *testing.T) {
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
						Type: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test",
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
						Type: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test",
					},
				},
			},
			err: context.Canceled,
		},
		"even fields unpack struct should be applied to sorted struct": {
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Type:  "test-1",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-2",
						Size:  16,
						Align: 16,
						Ptr:   4,
					},
					{
						Name:  "test",
						Type:  "test-3",
						Size:  24,
						Align: 16,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-4",
						Size:  4,
						Align: 8,
					},
					{
						Name:  "test",
						Type:  "test-5",
						Size:  20,
						Align: 20,
						Ptr:   1,
					},
					{
						Name:  "test",
						Type:  "test-6",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Type:  "test-4",
						Size:  4,
						Align: 8,
						Ptr:   0,
					},
					{
						Name:  "test",
						Type:  "test-5",
						Size:  20,
						Align: 20,
						Ptr:   1,
					},
					{
						Name:  "test",
						Type:  "test-6",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-2",
						Size:  16,
						Align: 16,
						Ptr:   4,
					},
					{
						Name:  "test",
						Type:  "test-1",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-3",
						Size:  24,
						Align: 16,
						Ptr:   8,
					},
				},
			},
		},
		"odd fields unpack struct should be applied to sorted struct": {
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Type:  "test-1",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-2",
						Size:  16,
						Align: 16,
						Ptr:   4,
					},
					{
						Name:  "test",
						Type:  "test-3",
						Size:  24,
						Align: 16,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-4",
						Size:  4,
						Align: 8,
					},
					{
						Name:  "test",
						Type:  "test-5",
						Size:  20,
						Align: 20,
						Ptr:   1,
					},
					{
						Name:  "test",
						Type:  "test-6",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-7",
						Size:  0,
						Align: 0,
						Ptr:   0,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Type:  "test-4",
						Size:  4,
						Align: 8,
						Ptr:   0,
					},
					{
						Name:  "test",
						Type:  "test-7",
						Size:  0,
						Align: 0,
						Ptr:   0,
					},
					{
						Name:  "test",
						Type:  "test-6",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-5",
						Size:  20,
						Align: 20,
						Ptr:   1,
					},
					{
						Name:  "test",
						Type:  "test-1",
						Size:  8,
						Align: 8,
						Ptr:   8,
					},
					{
						Name:  "test",
						Type:  "test-2",
						Size:  16,
						Align: 16,
						Ptr:   4,
					},
					{
						Name:  "test",
						Type:  "test-3",
						Size:  24,
						Align: 16,
						Ptr:   8,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := unpck.Apply(tcase.ctx, tcase.o)
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
