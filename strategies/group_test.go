package strategies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/tests/mocks"
)

func TestGroup(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		b   Builder
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to itself": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
		},
		"non empty struct without tag should be applied to itself": {
			b:   Builder{Curator: mocks.Maven{}},
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
					},
				},
			},
		},
		"non empty struct with irrelevant tag should be applied to itself": {
			b:   Builder{Curator: mocks.Maven{}},
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
						Tag:  `json:"test"`,
					},
				},
			},
		},
		"non empty struct with relevant tag should be applied to expected struct": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct with relevant tag should be applied to expected struct on canceled context": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate"`,
					},
				},
			},
			err: context.Canceled,
		},
		"non empty struct with relevant group tag should be applied to expected struct": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;doc_fields_annotate"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;doc_fields_annotate"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct with relevant group tag should be applied to expected struct with excess separators": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:";;group:def;doc_fields_annotate;;"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:";;group:def;doc_fields_annotate;;"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; - 🌺 gopium @1pkg"},
					},
				},
			},
		},
		"non empty struct with invalid tag should be applied to itself with error": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;doc_fields_annotate;test"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;doc_fields_annotate;test"`,
					},
				},
			},
			err: errors.New("tag \"gopium:\\\"group:def;doc_fields_annotate;test\\\"\" can't be parsed, neither as `default` nor named group"),
		},
		"non empty struct with invalid reverted tag should be applied to itself with error": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate;group:def"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate;group:def"`,
					},
				},
			},
			err: errors.New("tag \"gopium:\\\"doc_fields_annotate;group:def\\\"\" can't be parsed, named group `group:` anchor wasn't found"),
		},
		"mixed struct with inconsistent tags should be applied to itself with error": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"comment_fields_annotate"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"doc_fields_annotate"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"comment_fields_annotate"`,
					},
				},
			},
			err: errors.New(`inconsistent strategies list "comment_fields_annotate" for field "test" in group "default"`),
		},
		"mixed struct with inconsistent group tags should be applied to itself with error": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;doc_fields_annotate"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;comment_fields_annotate"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;doc_fields_annotate"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;comment_fields_annotate"`,
					},
				},
			},
			err: errors.New(`inconsistent strategies list "comment_fields_annotate" for field "test" in group "def"`),
		},
		"mixed struct with invalid strategies should be applied to itself with error": {
			b:   Builder{Curator: mocks.Maven{}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-1;doc_fields_annotate"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-2;doc_fields_annotate,test"`,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-1;doc_fields_annotate"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-2;doc_fields_annotate,test"`,
					},
				},
			},
			err: errors.New(`strategy "test" wasn't found`),
		},
		"mixed struct should be applied to expected struct accordingly to tags": {
			b:   Builder{Curator: mocks.Maven{SAlign: 12, SCache: []int64{24}}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name:    "test",
				Comment: []string{"// test"},
				Fields: []gopium.Field{
					{
						Name:  "test-1",
						Size:  12,
						Align: 8,
						Tag:   `gopium:"group:def-2;comment_fields_annotate,explicit_padings_system_alignment,cache_rounding_cpu_l1"`,
					},
					{
						Name:  "test-2",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-1;doc_fields_annotate,doc_struct_annotate"`,
					},
					{
						Name:  "test-3",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-2;comment_fields_annotate,explicit_padings_system_alignment,cache_rounding_cpu_l1"`,
					},
					{
						Name:  "test-4",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"comment_struct_stamp,cache_rounding_cpu_l1,explicit_padings_system_alignment"`,
					},
					{
						Name:  "test-5",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"comment_struct_stamp,cache_rounding_cpu_l1,explicit_padings_system_alignment"`,
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"// test", "// struct has been auto curated - 🌺 gopium @1pkg"},
				Fields: []gopium.Field{
					{
						Name:  "test-2",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-1;doc_fields_annotate,doc_struct_annotate"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test-1",
						Size:    12,
						Align:   8,
						Tag:     `gopium:"group:def-2;comment_fields_annotate,explicit_padings_system_alignment,cache_rounding_cpu_l1"`,
						Comment: []string{"// field size: 12 bytes; field align: 8 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test-3",
						Size:    8,
						Align:   4,
						Tag:     `gopium:"group:def-2;comment_fields_annotate,explicit_padings_system_alignment,cache_rounding_cpu_l1"`,
						Comment: []string{"// field size: 8 bytes; field align: 4 bytes; - 🌺 gopium @1pkg"},
					},
					gopium.PadField(4),
					{
						Name:  "test-4",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"comment_struct_stamp,cache_rounding_cpu_l1,explicit_padings_system_alignment"`,
					},
					gopium.PadField(4),
					{
						Name:  "test-5",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"comment_struct_stamp,cache_rounding_cpu_l1,explicit_padings_system_alignment"`,
					},
					gopium.PadField(4),
					gopium.PadField(8),
					gopium.PadField(4),
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			grp := ptgrp.Builder(tcase.b)
			// exec
			r, err := grp.Apply(tcase.ctx, tcase.o)
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
