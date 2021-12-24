package strategies

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests/mocks"
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
						Tag:   `gopium:"fields_annotate_doc"`,
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
						Tag:   `gopium:"fields_annotate_doc"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
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
						Tag:   `gopium:"fields_annotate_doc"`,
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
						Tag:   `gopium:"fields_annotate_doc"`,
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
						Tag:   `gopium:"group:def;fields_annotate_doc"`,
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
						Tag:   `gopium:"group:def;fields_annotate_doc"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
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
						Tag:   `gopium:";;group:def;fields_annotate_doc;;"`,
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
						Tag:   `gopium:";;group:def;fields_annotate_doc;;"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
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
						Tag:   `gopium:"group:def;fields_annotate_doc;test"`,
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
						Tag:   `gopium:"group:def;fields_annotate_doc;test"`,
					},
				},
			},
			err: errors.New("tag \"gopium:\\\"group:def;fields_annotate_doc;test\\\"\" can't be parsed, neither as `default` nor named group"),
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
						Tag:   `gopium:"fields_annotate_doc;group:def"`,
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
						Tag:   `gopium:"fields_annotate_doc;group:def"`,
					},
				},
			},
			err: errors.New("tag \"gopium:\\\"fields_annotate_doc;group:def\\\"\" can't be parsed, named group `group:` anchor wasn't found"),
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
						Tag:   `gopium:"fields_annotate_doc"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"fields_annotate_comment"`,
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
						Tag:   `gopium:"fields_annotate_doc"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"fields_annotate_comment"`,
					},
				},
			},
			err: errors.New(`inconsistent strategies list "fields_annotate_comment" for field "test" in default group`),
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
						Tag:   `gopium:"group:def;fields_annotate_doc"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;fields_annotate_comment"`,
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
						Tag:   `gopium:"group:def;fields_annotate_doc"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def;fields_annotate_comment"`,
					},
				},
			},
			err: errors.New(`inconsistent strategies list "fields_annotate_comment" for field "test" in group "def"`),
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
						Tag:   `gopium:"group:def-1;fields_annotate_doc"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-2;fields_annotate_doc,test"`,
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
						Tag:   `gopium:"group:def-1;fields_annotate_doc"`,
					},
					{
						Name:  "test",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-2;fields_annotate_doc,test"`,
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
						Ptr:   4,
						Tag:   `gopium:"group:def-2;fields_annotate_comment,explicit_paddings_system_alignment,cache_rounding_cpu_l1_discrete"`,
					},
					{
						Name:  "test-2",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-1;fields_annotate_doc,struct_annotate_doc"`,
					},
					{
						Name:  "test-3",
						Size:  8,
						Align: 4,
						Ptr:   4,
						Tag:   `gopium:"group:def-2;fields_annotate_comment,explicit_paddings_system_alignment,cache_rounding_cpu_l1_discrete"`,
					},
					{
						Name:  "test-4",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"cache_rounding_cpu_l1_discrete,explicit_paddings_system_alignment"`,
					},
					{
						Name:  "test-5",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"cache_rounding_cpu_l1_discrete,explicit_paddings_system_alignment"`,
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Comment: []string{"// test"},
				Fields: []gopium.Field{
					{
						Name:  "test-4",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"cache_rounding_cpu_l1_discrete,explicit_paddings_system_alignment"`,
					},
					collections.PadField(4),
					{
						Name:  "test-5",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"cache_rounding_cpu_l1_discrete,explicit_paddings_system_alignment"`,
					},
					collections.PadField(4),
					collections.PadField(8),
					collections.PadField(4),
					{
						Name:  "test-2",
						Size:  8,
						Align: 4,
						Tag:   `gopium:"group:def-1;fields_annotate_doc,struct_annotate_doc"`,
						Doc:   []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 0 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test-1",
						Size:    12,
						Align:   8,
						Ptr:     4,
						Tag:     `gopium:"group:def-2;fields_annotate_comment,explicit_paddings_system_alignment,cache_rounding_cpu_l1_discrete"`,
						Comment: []string{"// field size: 12 bytes; field align: 8 bytes; field ptr: 4 bytes; - 🌺 gopium @1pkg"},
					},
					{
						Name:    "test-3",
						Size:    8,
						Align:   4,
						Ptr:     4,
						Tag:     `gopium:"group:def-2;fields_annotate_comment,explicit_paddings_system_alignment,cache_rounding_cpu_l1_discrete"`,
						Comment: []string{"// field size: 8 bytes; field align: 4 bytes; field ptr: 4 bytes; - 🌺 gopium @1pkg"},
					},
					collections.PadField(4),
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			grp := ptag.Builder(tcase.b)
			// exec
			r, err := grp.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
