package collections

import (
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestFlatSorted(t *testing.T) {
	// prepare
	table := map[string]struct {
		f Flat
		r []gopium.Struct
	}{
		"nil flat collection should return empty sorted": {
			f: nil,
			r: []gopium.Struct{},
		},
		"empty flat collection should return empty sorted": {
			f: Flat{},
			r: []gopium.Struct{},
		},
		"single item flat collection should return single item sorted": {
			f: Flat{"1-test": gopium.Struct{Name: "test1"}},
			r: []gopium.Struct{gopium.Struct{Name: "test1"}},
		},
		"multiple presorted items flat collection should return multiple items sorted": {
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
				"3-test": gopium.Struct{Name: "test3"},
			},
			r: []gopium.Struct{
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
			},
		},
		"multiple reverted items flat collection should return multiple items sorted": {
			f: Flat{
				"3-test": gopium.Struct{Name: "test3"},
				"2-test": gopium.Struct{Name: "test2"},
				"1-test": gopium.Struct{Name: "test1"},
			},
			r: []gopium.Struct{
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
			},
		},
		"multiple mixed items flat collection should return multiple items sorted": {
			f: Flat{
				"99-test":   gopium.Struct{Name: "test99"},
				"5-test":    gopium.Struct{Name: "test5"},
				"1000-test": gopium.Struct{Name: "test1000"},
				"3-test":    gopium.Struct{Name: "test3"},
				"1-test":    gopium.Struct{Name: "test1"},
				"2-test":    gopium.Struct{Name: "test2"},
				"4-test":    gopium.Struct{Name: "test4"},
				"0-test":    gopium.Struct{Name: "test0"},
			},
			r: []gopium.Struct{
				gopium.Struct{Name: "test0"},
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
				gopium.Struct{Name: "test4"},
				gopium.Struct{Name: "test5"},
				gopium.Struct{Name: "test99"},
				gopium.Struct{Name: "test1000"},
			},
		},
		"multiple non pattern ids items flat collection should return items sorted naturally": {
			f: Flat{
				"a-test": gopium.Struct{Name: "testa"},
				"b-test": gopium.Struct{Name: "testb"},
				"c-test": gopium.Struct{Name: "testc"},
			},
			r: []gopium.Struct{
				gopium.Struct{Name: "testa"},
				gopium.Struct{Name: "testb"},
				gopium.Struct{Name: "testc"},
			},
		},
		"multiple mixed non pattern ids items flat collection should return items sorted naturally": {
			f: Flat{
				"3-test": gopium.Struct{Name: "test3"},
				"2-test": gopium.Struct{Name: "test2"},
				"1-test": gopium.Struct{Name: "test1"},
				"a-test": gopium.Struct{Name: "testa"},
				"b-test": gopium.Struct{Name: "testb"},
				"c-test": gopium.Struct{Name: "testc"},
			},
			r: []gopium.Struct{
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
				gopium.Struct{Name: "testa"},
				gopium.Struct{Name: "testb"},
				gopium.Struct{Name: "testc"},
			},
		},
		"cmplex multiple mixed non pattern ids items flat collection should return items sorted naturally": {
			f: Flat{
				"z-test":    gopium.Struct{Name: "testz"},
				"3-test":    gopium.Struct{Name: "test3"},
				"2-test":    gopium.Struct{Name: "test2"},
				"1-test":    gopium.Struct{Name: "test1"},
				"a-test":    gopium.Struct{Name: "testa"},
				"b-test":    gopium.Struct{Name: "testb"},
				"c-test":    gopium.Struct{Name: "testc"},
				"99-test":   gopium.Struct{Name: "test99"},
				"5-test":    gopium.Struct{Name: "test5"},
				"1000-test": gopium.Struct{Name: "test1000"},
				"4-test":    gopium.Struct{Name: "test4"},
				"0-test":    gopium.Struct{Name: "test0"},
				"xytest":    gopium.Struct{Name: "testxy"},
			},
			r: []gopium.Struct{
				gopium.Struct{Name: "test0"},
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
				gopium.Struct{Name: "test4"},
				gopium.Struct{Name: "test5"},
				gopium.Struct{Name: "test99"},
				gopium.Struct{Name: "test1000"},
				gopium.Struct{Name: "testa"},
				gopium.Struct{Name: "testb"},
				gopium.Struct{Name: "testc"},
				gopium.Struct{Name: "testxy"},
				gopium.Struct{Name: "testz"},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r := tcase.f.Sorted()
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
		})
	}
}
