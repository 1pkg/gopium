package strategies

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/1pkg/gopium/gopium"
)

func TestFilter(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		filter filter
		ctx    context.Context
		o      gopium.Struct
		r      gopium.Struct
		err    error
	}{
		"empty struct should be applied to empty struct with empty filter": {
			filter: filter{},
			ctx:    context.Background(),
		},
		"non empty struct should be applied to itself with empty filter": {
			filter: filter{},
			ctx:    context.Background(),
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
			filter: filter{},
			ctx:    cctx,
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
		"non empty struct should be applied accordingly to filter name": {
			filter: filter{nregex: regexp.MustCompile(`^test-2$`)},
			ctx:    context.Background(),
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
			filter: filter{tregex: regexp.MustCompile(`^test-2$`)},
			ctx:    context.Background(),
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
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.filter.Apply(tcase.ctx, tcase.o)
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
