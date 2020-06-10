package fmtio

import (
	"reflect"
	"strings"
	"testing"

	"1pkg/gopium/collections"
	"1pkg/gopium/gopium"
)

func TestDiff(t *testing.T) {
	// prepare
	oh := collections.NewHierarchic("")
	rh := collections.NewHierarchic("")
	rhb := collections.NewHierarchic("")
	oh.Push("test", "test", gopium.Struct{
		Name: "test",
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
				Align: 8,
			},
			{
				Name:  "test3",
				Size:  3,
				Align: 1,
			},
		},
	})
	rh.Push("test", "test", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test1",
				Size:  3,
				Align: 1,
			},
			{
				Name:  "test3",
				Size:  3,
				Align: 1,
			},
		},
	})
	rhb.Push("test", "test", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
		},
	})
	rhb.Push("test1", "test1", gopium.Struct{
		Name: "test",
		Fields: []gopium.Field{
			{
				Name:  "test2",
				Type:  "float64",
				Size:  8,
				Align: 8,
			},
		},
	})
	table := map[string]struct {
		fmt gopium.Diff
		o   gopium.Categorized
		r   gopium.Categorized
		b   []byte
		err error
	}{
		"size align md table should return expected result for empty collections": {
			fmt: SizeAlignMdt,
			o:   collections.NewHierarchic(""),
			r:   collections.NewHierarchic(""),
			b: []byte(`
| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |
| :---: | :---: | :---: | :---: | :---: |
`),
		},
		"md table should return expected result for non empty collections": {
			fmt: SizeAlignMdt,
			o:   oh,
			r:   rh,
			b: []byte(`
| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |
| :---: | :---: | :---: | :---: | :---: |
| test | 24 bytes | 16 bytes | -8 bytes | -33.33% |
| Total | 24 bytes | 16 bytes | -8 bytes | -33.33% |
`),
		},
		"md table should return expected result for non empty overlapping collections": {
			fmt: SizeAlignMdt,
			o:   oh,
			r:   rhb,
			b: []byte(`
| Struct Name | Original Size with Pad | Current Size with Pad | Absolute Difference | Relative Difference |
| :---: | :---: | :---: | :---: | :---: |
| test | 24 bytes | 32 bytes | +8 bytes | +33.33% |
| Total | 24 bytes | 32 bytes | +8 bytes | +33.33% |
`),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			b, err := tcase.fmt(tcase.o, tcase.r)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// format actual and expected identically
			actual := strings.Trim(string(b), "\n")
			expected := strings.Trim(string(tcase.b), "\n")
			if err == nil && !reflect.DeepEqual(actual, expected) {
				t.Errorf("actual %v doesn't equal to expected %v", actual, expected)
			}
		})
	}
}
