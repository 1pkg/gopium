package strategies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests/mocks"
)

func TestPipe(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	terr := errors.New("test error")
	table := map[string]struct {
		pipe pipe
		ctx  context.Context
		o    gopium.Struct
		r    gopium.Struct
		err  error
	}{
		"empty struct should be applied to empty struct with empty pipe": {
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself with empty pipe": {
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
		"non empty struct should be applied to itself with empty pipe on canceled context": {
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
		"non empty struct should be applied accordingly to pipe": {
			pipe: pipe([]gopium.Strategy{fnotecom, fnotedoc}),
			ctx:  context.Background(),
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
						Name:    "test",
						Type:    "test",
						Doc:     []string{"// field size: 0 bytes; field align: 0 bytes; - 🌺 gopium @1pkg"},
						Comment: []string{"// field size: 0 bytes; field align: 0 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct should be applied to itself on canceled context": {
			pipe: pipe([]gopium.Strategy{fnotecom, fnotedoc}),
			ctx:  cctx,
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
		"non empty struct should be applied to expected result on pipe error": {
			pipe: pipe([]gopium.Strategy{
				fnotecom,
				&mocks.Strategy{Err: terr},
				fnotedoc,
			}),
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
						Name:    "test",
						Type:    "test",
						Comment: []string{"// field size: 0 bytes; field align: 0 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
			err: terr,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.pipe.Apply(tcase.ctx, tcase.o)
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
