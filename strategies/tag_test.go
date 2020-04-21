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
			tag: tag{},
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself with relevant tag": {
			tag: tag{tag: "test"},
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
		"non empty struct should be applied to itself with relevant doc on canceled context": {
			tag: tag{tag: "test"},
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
			err: cctx.Err(),
		},
		"non empty struct should be applied to itself valid tag shouldn't be overwritten": {
			tag: tag{tag: "test"},
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
		"non empty struct should be applied to itself new tag shouldn be appended": {
			tag: tag{tag: "test"},
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
		"non empty struct should be applied to itself with relevant tag should be overwritten": {
			tag: tag{tag: "test", force: true},
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
		"complex struct should be applied to itself with relevant tags": {
			tag: tag{tag: "test", force: true},
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
		},
		"complex struct should be applied to itself with relevant tags with grop prefix": {
			tag: tag{
				tag:   "test",
				group: "group",
				force: true,
			},
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
						Tag:  `json:"test" gopium:"group:group;test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"group:group;test"`,
					},
					{
						Name: "test2",
						Type: "float64",
						Tag:  `gopium:"group:group;test"`,
					},
				},
			},
		},
		"complex struct should be applied to itself with relevant tags on discrete": {
			tag: tag{
				tag:      "test",
				force:    true,
				discrete: true,
			},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
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
						Tag:  `gopium:"group:default-1;test"`,
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
		"complex struct should be applied to itself with relevant tags on discrete with prefix group": {
			tag: tag{
				tag:      "test",
				group:    "group",
				force:    true,
				discrete: true,
			},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Type: "int",
						Size: 8,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
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
						Tag:  `gopium:"group:group-1;test"`,
					},
					{
						Name: "test2",
						Type: "string",
						Doc:  []string{"test"},
						Tag:  `gopium:"group:group-2;test"`,
					},
					{
						Name: "test2",
						Type: "float64",
						Tag:  `gopium:"group:group-3;test"`,
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
