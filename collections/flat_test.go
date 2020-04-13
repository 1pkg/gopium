package collections

import (
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestFlatSorted(t *testing.T) {
	// prepare
	table := map[string]struct {
		input  Flat
		output []gopium.Struct
	}{
		"nil flat collection should return empty sorted": {
			input:  nil,
			output: []gopium.Struct{},
		},
		"empty flat collection should return empty sorted": {
			input:  Flat{},
			output: []gopium.Struct{},
		},
		"single item flat collection should return single item sorted": {
			input:  Flat{"1-test": gopium.Struct{Name: "test1"}},
			output: []gopium.Struct{gopium.Struct{Name: "test1"}},
		},
		"multiple presorted items flat collection should return multiple items sorted": {
			input: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
				"3-test": gopium.Struct{Name: "test3"},
			},
			output: []gopium.Struct{
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
			},
		},
		"multiple reverted items flat collection should return multiple items sorted": {
			input: Flat{
				"3-test": gopium.Struct{Name: "test3"},
				"2-test": gopium.Struct{Name: "test2"},
				"1-test": gopium.Struct{Name: "test1"},
			},
			output: []gopium.Struct{
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
			},
		},
		"multiple mixed items flat collection should return multiple items sorted": {
			input: Flat{
				"99-test":   gopium.Struct{Name: "test99"},
				"5-test":    gopium.Struct{Name: "test5"},
				"1000-test": gopium.Struct{Name: "test1000"},
				"3-test":    gopium.Struct{Name: "test3"},
				"1-test":    gopium.Struct{Name: "test1"},
				"2-test":    gopium.Struct{Name: "test2"},
				"4-test":    gopium.Struct{Name: "test4"},
				"0-test":    gopium.Struct{Name: "test0"},
			},
			output: []gopium.Struct{
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
		"multiple non pattern ids items flat collection should still return items sorted naturally": {
			input: Flat{
				"a-test": gopium.Struct{Name: "testa"},
				"b-test": gopium.Struct{Name: "testb"},
				"c-test": gopium.Struct{Name: "testc"},
			},
			output: []gopium.Struct{
				gopium.Struct{Name: "testa"},
				gopium.Struct{Name: "testb"},
				gopium.Struct{Name: "testc"},
			},
		},
		"multiple mixed non pattern ids items flat collection should still return items sorted naturally": {
			input: Flat{
				"3-test": gopium.Struct{Name: "test3"},
				"2-test": gopium.Struct{Name: "test2"},
				"1-test": gopium.Struct{Name: "test1"},
				"a-test": gopium.Struct{Name: "testa"},
				"b-test": gopium.Struct{Name: "testb"},
				"c-test": gopium.Struct{Name: "testc"},
			},
			output: []gopium.Struct{
				gopium.Struct{Name: "test1"},
				gopium.Struct{Name: "test2"},
				gopium.Struct{Name: "test3"},
				gopium.Struct{Name: "testa"},
				gopium.Struct{Name: "testb"},
				gopium.Struct{Name: "testc"},
			},
		},
		"cmplex multiple mixed non pattern ids items flat collection should still return items sorted naturally": {
			input: Flat{
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
			output: []gopium.Struct{
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
			output := tcase.input.Sorted()
			// check
			if !reflect.DeepEqual(output, tcase.output) {
				t.Errorf("actual %v not equals expected %v", output, tcase.output)
			}
		})
	}
}
