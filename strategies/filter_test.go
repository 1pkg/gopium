package strategies

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"1pkg/gopium"
)

func TestFilter(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		f   filter
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct with empty filter": {
			f:   filter{},
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself with empty filter": {
			f:   filter{},
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
		"non empty struct should be applied to itself on canceled context with empty filter": {
			f:   filter{},
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
			err: cctx.Err(),
		},
		"non empty struct should be applied accordingly to filter name": {
			f:   filter{nregex: regexp.MustCompile(`^test-2$`)},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
		},
		"non empty struct should be applied accordingly to filter type": {
			f:   filter{tregex: regexp.MustCompile(`^test-2$`)},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
		},
		"non empty struct should be applied accordingly to filter exported true": {
			f:   filter{exp: &tvar},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
		},
		"non empty struct should be applied accordingly to filter exported false": {
			f:   filter{exp: &fvar},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
				},
			},
		},
		"non empty struct should be applied accordingly to filter embedded true": {
			f:   filter{emb: &tvar},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
				},
			},
		},
		"non empty struct should be applied accordingly to filter embedded false": {
			f:   filter{emb: &fvar},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
		},
		"non empty struct should be applied accordingly to filter": {
			f:   filter{nregex: regexp.MustCompile(`^test-2$`), exp: &fvar},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
					{
						Name:     "test-2",
						Type:     "test-2",
						Exported: true,
					},
					{
						Name:     "test-3",
						Type:     "test-3",
						Embedded: true,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test-1",
						Type:     "test-1",
						Embedded: true,
						Exported: true,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.f.Apply(tcase.ctx, tcase.o)
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
