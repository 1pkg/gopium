package collections

import (
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
)

func TestSizeAlign(t *testing.T) {
	// prepare
	table := map[string]struct {
		st    gopium.Struct
		size  int64
		align int64
	}{
		"empty struct should return expected size and align": {
			size:  0,
			align: 1,
		},
		"non empty struct should return expected size and align": {
			st: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Type:  "int",
						Size:  8,
						Align: 4,
					},
					{
						Name:    "test2",
						Type:    "string",
						Comment: []string{"test"},
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 8,
					},
				},
			},
			size:  16,
			align: 8,
		},
		"struct with pads should return expected size and align": {
			st: gopium.Struct{
				Name:    "test",
				Comment: []string{"test"},
				Fields: []gopium.Field{
					{
						Name:  "test1",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test2",
						Type:  "float64",
						Size:  8,
						Align: 6,
					},
					{
						Name:  "test3",
						Size:  3,
						Align: 1,
					},
				},
			},
			size:  18,
			align: 6,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			size, align := SizeAlign(tcase.st)
			// check
			if !reflect.DeepEqual(size, tcase.size) {
				t.Errorf("actual %v doesn't equal to %v", size, tcase.size)
			}
			if !reflect.DeepEqual(align, tcase.align) {
				t.Errorf("actual %v doesn't equal to %v", size, tcase.size)
			}
		})
	}
}

func TestPadField(t *testing.T) {
	// prepare
	table := map[string]struct {
		pad int64
		f   gopium.Field
	}{
		"empty pad should return empty field pad": {
			f: gopium.Field{
				Name:  "_",
				Type:  "[0]byte",
				Size:  0,
				Align: 1,
			},
		},
		"positive pad should return valid field pad": {
			pad: 10,
			f: gopium.Field{
				Name:  "_",
				Type:  "[10]byte",
				Size:  10,
				Align: 1,
			},
		},
		"negative pad should return empty field": {
			pad: -10,
			f: gopium.Field{
				Name:  "_",
				Type:  "[0]byte",
				Size:  0,
				Align: 1,
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			f := PadField(tcase.pad)
			// check
			if !reflect.DeepEqual(f, tcase.f) {
				t.Errorf("actual %v doesn't equal to %v", f, tcase.f)
			}
		})
	}
}
