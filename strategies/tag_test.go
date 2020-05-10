package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestTag(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		tag tag
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to itself": {
			tag: tags,
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself with expected tag": {
			tag: tags.Names(gopium.StrategyName("test")),
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
						Tag:  `gopium:"test"`,
					},
				},
			},
		},
		"non empty struct should be applied to itself with expected tag on canceled context": {
			tag: tags.Names(gopium.StrategyName("test")),
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
						Tag:  `gopium:"test"`,
					},
				},
			},
			err: context.Canceled,
		},
		"non empty struct should be applied to itself valid tag shouldn't be overwritten": {
			tag: tags.Names(gopium.StrategyName("test")),
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Tag:  `gopium:"tag"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Tag:  `gopium:"tag"`,
					},
				},
			},
		},
		"non empty struct should be applied to itself new tag should be appended": {
			tag: tags.Names(gopium.StrategyName("test")),
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Tag:  `json:"test"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Tag:  `json:"test" gopium:"test"`,
					},
				},
			},
		},
		"non empty struct should be applied to itself with expected tag should be overwritten": {
			tag: tagf.Names(gopium.StrategyName("test")),
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Tag:  `json:"test" gopium:"tag"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Tag:  `json:"test" gopium:"test"`,
					},
				},
			},
		},
		"complex struct should be applied to itself with expected tag on force": {
			tag: tagf.Names(gopium.StrategyName("test")),
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
						Tag:  `json:"test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"string"`,
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
						Tag:  `json:"test" gopium:"test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"test"`,
					},
					{
						Name: "test2",
						Type: "float64",
						Tag:  `gopium:"test"`,
					},
				},
			},
		}, "complex struct should be applied to itself with expected tag on soft discrete": {
			tag: tagsd.Names(gopium.StrategyName("test")),
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
						Tag:  `json:"test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"string"`,
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
						Tag:  `json:"test" gopium:"group:default-1;test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"string"`,
					},
					{
						Name: "test2",
						Type: "float64",
						Tag:  `gopium:"group:default-3;test"`,
					},
				},
			},
		},
		"complex struct should be applied to itself with expected tag on force discrete": {
			tag: tagfd.Names(gopium.StrategyName("test")),
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
						Tag:  `json:"test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"string"`,
					},
					{
						Name: "test2",
						Type: "float64",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
						Tag:  `json:"test" gopium:"group:default-1;test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"group:default-2;test"`,
					},
					{
						Name: "test2",
						Type: "float64",
						Tag:  `gopium:"group:default-3;test"`,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.tag.Apply(tcase.ctx, tcase.o)
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
