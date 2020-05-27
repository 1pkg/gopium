package collections

import (
	"reflect"
	"testing"

	"1pkg/gopium"
)

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
