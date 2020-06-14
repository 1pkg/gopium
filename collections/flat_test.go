package collections

import (
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
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
			f: Flat{"test:1": gopium.Struct{Name: "test1"}},
			r: []gopium.Struct{{Name: "test1"}},
		},
		"multiple presorted items flat collection should return multiple items sorted": {
			f: Flat{
				"test:1": gopium.Struct{Name: "test1"},
				"test:2": gopium.Struct{Name: "test2"},
				"test:3": gopium.Struct{Name: "test3"},
			},
			r: []gopium.Struct{
				{Name: "test1"},
				{Name: "test2"},
				{Name: "test3"},
			},
		},
		"multiple reverted items flat collection should return multiple items sorted": {
			f: Flat{
				"test:3": gopium.Struct{Name: "test3"},
				"test:2": gopium.Struct{Name: "test2"},
				"test:1": gopium.Struct{Name: "test1"},
			},
			r: []gopium.Struct{
				{Name: "test1"},
				{Name: "test2"},
				{Name: "test3"},
			},
		},
		"multiple mixed items flat collection should return multiple items sorted": {
			f: Flat{
				"test:99":   gopium.Struct{Name: "test99"},
				"test:5":    gopium.Struct{Name: "test5"},
				"test:1000": gopium.Struct{Name: "test1000"},
				"test:3":    gopium.Struct{Name: "test3"},
				"test:1":    gopium.Struct{Name: "test1"},
				"test:2":    gopium.Struct{Name: "test2"},
				"test:4":    gopium.Struct{Name: "test4"},
				"test:0":    gopium.Struct{Name: "test0"},
			},
			r: []gopium.Struct{
				{Name: "test0"},
				{Name: "test1"},
				{Name: "test2"},
				{Name: "test3"},
				{Name: "test4"},
				{Name: "test5"},
				{Name: "test99"},
				{Name: "test1000"},
			},
		},
		"multiple non pattern ids items flat collection should return items sorted naturally": {
			f: Flat{
				"test:a": gopium.Struct{Name: "testa"},
				"test:b": gopium.Struct{Name: "testb"},
				"test:c": gopium.Struct{Name: "testc"},
			},
			r: []gopium.Struct{
				{Name: "testa"},
				{Name: "testb"},
				{Name: "testc"},
			},
		},
		"multiple mixed non pattern ids items flat collection should return items sorted naturally": {
			f: Flat{
				"test:3": gopium.Struct{Name: "test3"},
				"test:2": gopium.Struct{Name: "test2"},
				"test:1": gopium.Struct{Name: "test1"},
				"test:a": gopium.Struct{Name: "testa"},
				"test:b": gopium.Struct{Name: "testb"},
				"test:c": gopium.Struct{Name: "testc"},
			},
			r: []gopium.Struct{
				{Name: "test1"},
				{Name: "test2"},
				{Name: "test3"},
				{Name: "testa"},
				{Name: "testb"},
				{Name: "testc"},
			},
		},
		"complex multiple mixed non pattern ids items flat collection should return items sorted naturally": {
			f: Flat{
				"test:z":    gopium.Struct{Name: "testz"},
				"test:3":    gopium.Struct{Name: "test3"},
				"test:2":    gopium.Struct{Name: "test2"},
				"test:1":    gopium.Struct{Name: "test1"},
				"test:a":    gopium.Struct{Name: "testa"},
				"test:b":    gopium.Struct{Name: "testb"},
				"test:c":    gopium.Struct{Name: "testc"},
				"test:99":   gopium.Struct{Name: "test99"},
				"test:5":    gopium.Struct{Name: "test5"},
				"test:1000": gopium.Struct{Name: "test1000"},
				"test:4":    gopium.Struct{Name: "test4"},
				"test:0":    gopium.Struct{Name: "test0"},
				"test:xy":   gopium.Struct{Name: "testxy"},
			},
			r: []gopium.Struct{
				{Name: "test0"},
				{Name: "test1"},
				{Name: "test2"},
				{Name: "test3"},
				{Name: "test4"},
				{Name: "test5"},
				{Name: "test99"},
				{Name: "test1000"},
				{Name: "testa"},
				{Name: "testb"},
				{Name: "testc"},
				{Name: "testxy"},
				{Name: "testz"},
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
