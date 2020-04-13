package collections

import (
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestHierarchicPushCatMixed(t *testing.T) {
	h := Hierarchic{}
	// check non existed cat
	cat, ok := h.Cat("cat")
	expected := Flat(nil)
	if ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
	// create new and check existed cat
	h.Push("test-1", "cat", gopium.Struct{Name: "test-1"})
	cat, ok = h.Cat("cat")
	expected = Flat{"test-1": gopium.Struct{Name: "test-1"}}
	if !ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
	// check another non existed cat
	cat, ok = h.Cat("cat-1")
	expected = Flat(nil)
	if ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
	// update and check existed cat
	h.Push("test-2", "cat", gopium.Struct{Name: "test-2"})
	h.Push("test-3", "cat", gopium.Struct{Name: "test-3"})
	cat, ok = h.Cat("cat")
	expected = Flat{
		"test-1": gopium.Struct{Name: "test-1"},
		"test-2": gopium.Struct{Name: "test-2"},
		"test-3": gopium.Struct{Name: "test-3"},
	}
	if !ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
	// check another non existed cat
	cat, ok = h.Cat("cat-1")
	expected = Flat(nil)
	if ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
	// create new and update and check existed cat
	h.Push("test-1", "cat", gopium.Struct{Name: "test-5"})
	h.Push("test-1", "cat-1", gopium.Struct{Name: "test-1"})
	cat, ok = h.Cat("cat")
	expected = Flat{
		"test-1": gopium.Struct{Name: "test-5"},
		"test-2": gopium.Struct{Name: "test-2"},
		"test-3": gopium.Struct{Name: "test-3"},
	}
	if !ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
	// check another existed cat
	cat, ok = h.Cat("cat-1")
	expected = Flat{"test-1": gopium.Struct{Name: "test-1"}}
	if !ok || !reflect.DeepEqual(cat, expected) {
		t.Errorf("actual %v doesn't equal to %v", cat, expected)
	}
}

func TestHierarchicFlat(t *testing.T) {
	// prepare
	table := map[string]struct {
		input  Hierarchic
		output Flat
	}{
		"nil hierarchic collection should return empty flat collection": {
			input:  nil,
			output: Flat{},
		},
		"empty hierarchic collection should return empty flat collection": {
			input:  Hierarchic{},
			output: Flat{},
		},
		"single loc single item hierarchic collection should return single item flat collection": {
			input:  Hierarchic{"loc": {"1-test": gopium.Struct{Name: "test1"}}},
			output: Flat{"1-test": gopium.Struct{Name: "test1"}},
		},
		"single loc multiple items hierarchic collection should return single item flat collection": {
			input: Hierarchic{"loc": {
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			}},
			output: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			},
		},
		"multiple locs single item hierarchic collection should return multiple items flat collection": {
			input: Hierarchic{
				"loc-1": {"1-test": gopium.Struct{Name: "test1"}},
				"loc-2": {"2-test": gopium.Struct{Name: "test2"}},
			},
			output: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			},
		},
		"multiple locs multiple items hierarchic collection should return multiple items flat collection": {
			input: Hierarchic{
				"loc-1": {
					"1-test": gopium.Struct{Name: "test1"},
					"2-test": gopium.Struct{Name: "test2"},
				},
				"loc-2": {
					"3-test": gopium.Struct{Name: "test3"},
					"4-test": gopium.Struct{Name: "test4"},
				},
			},
			output: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
				"3-test": gopium.Struct{Name: "test3"},
				"4-test": gopium.Struct{Name: "test4"},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			output := tcase.input.Flat()
			// check
			if !reflect.DeepEqual(output, tcase.output) {
				t.Errorf("actual %v doesn't equal to %v", output, tcase.output)
			}
		})
	}
}
