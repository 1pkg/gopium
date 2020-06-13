package collections

import (
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
)

func TestCopyField(t *testing.T) {
	// prepare
	table := map[string]struct {
		o gopium.Field
		r gopium.Field
	}{
		"empty field should be copied to empty field": {
			o: gopium.Field{},
			r: gopium.Field{},
		},
		"non empty field should be copied to same field": {
			o: gopium.Field{
				Name:     "test",
				Type:     "type",
				Exported: true,
			},
			r: gopium.Field{
				Name:     "test",
				Type:     "type",
				Exported: true,
			},
		},
		"non empty field with notes should be copied to same field": {
			o: gopium.Field{
				Name:     "test",
				Type:     "type",
				Exported: true,
				Doc:      []string{"test-doc"},
				Comment:  []string{"test-com-1", "test-com-2"},
			},
			r: gopium.Field{
				Name:     "test",
				Type:     "type",
				Exported: true,
				Doc:      []string{"test-doc"},
				Comment:  []string{"test-com-1", "test-com-2"},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r := CopyField(tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to %v", r, tcase.r)
			}
		})
	}
}

func TestCopyStruct(t *testing.T) {
	// prepare
	table := map[string]struct {
		o gopium.Struct
		r gopium.Struct
	}{
		"empty struct should be copied to empty struct": {
			o: gopium.Struct{},
			r: gopium.Struct{},
		},
		"non empty struct should be copied to same struct": {
			o: gopium.Struct{
				Name: "test",
			},
			r: gopium.Struct{
				Name: "test",
			},
		},
		"non empty struct with notes should be copied to same struct": {
			o: gopium.Struct{
				Name:    "test",
				Doc:     []string{"test-doc"},
				Comment: []string{"test-com-1", "test-com-2"},
			},
			r: gopium.Struct{
				Name:    "test",
				Doc:     []string{"test-doc"},
				Comment: []string{"test-com-1", "test-com-2"},
			},
		},
		"non empty struct with notes and fields should be copied to same struct": {
			o: gopium.Struct{
				Name:    "test",
				Doc:     []string{"test-doc"},
				Comment: []string{"test-com-1", "test-com-2"},
				Fields: []gopium.Field{
					{
						Name:     "test",
						Type:     "type",
						Exported: true,
					},
					{
						Name:     "test",
						Type:     "type",
						Embedded: true,
						Doc:      []string{"test-doc"},
						Comment:  []string{"test-com-1", "test-com-2"},
					},
				},
			},
			r: gopium.Struct{
				Name:    "test",
				Doc:     []string{"test-doc"},
				Comment: []string{"test-com-1", "test-com-2"},
				Fields: []gopium.Field{
					{
						Name:     "test",
						Type:     "type",
						Exported: true,
					},
					{
						Name:     "test",
						Type:     "type",
						Embedded: true,
						Doc:      []string{"test-doc"},
						Comment:  []string{"test-com-1", "test-com-2"},
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r := CopyStruct(tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to %v", r, tcase.r)
			}
		})
	}
}
