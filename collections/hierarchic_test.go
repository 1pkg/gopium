package collections

import (
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestHierarchicPushCatMixed(t *testing.T) {
	// prepare
	h := Hierarchic{}
	h.Push("test-1", "test-1", gopium.Struct{Name: "test-1"})
	h.Push("test-2", "test-2", gopium.Struct{Name: "test-2"})
	h.Push("test-3", "test-2", gopium.Struct{Name: "test-3"})
	table := map[string]struct {
		cat string
		f   Flat
		ok  bool
	}{
		"invalid cat should return empty flat collection": {
			cat: "cat",
			f:   Flat(nil),
		},
		"test-1 cat should return expected flat collection": {
			cat: "test-1",
			f:   Flat{"test-1": gopium.Struct{Name: "test-1"}},
			ok:  true,
		},
		"test-2 cat should return expected flat collection": {
			cat: "test-2",
			f: Flat{
				"test-2": gopium.Struct{Name: "test-2"},
				"test-3": gopium.Struct{Name: "test-3"},
			},
			ok: true,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			f, ok := h.Cat(tcase.cat)
			// check
			if !reflect.DeepEqual(ok, tcase.ok) {
				t.Errorf("actual %v doesn't equal to %v", ok, tcase.ok)
			}
			if !reflect.DeepEqual(f, tcase.f) {
				t.Errorf("actual %v doesn't equal to %v", f, tcase.f)
			}
		})
	}
}

func TestHierarchicFlat(t *testing.T) {
	// prepare
	table := map[string]struct {
		h Hierarchic
		f Flat
	}{
		"nil hierarchic collection should return empty flat collection": {
			h: nil,
			f: Flat{},
		},
		"empty hierarchic collection should return empty flat collection": {
			h: Hierarchic{},
			f: Flat{},
		},
		"single loc single item hierarchic collection should return single item flat collection": {
			h: Hierarchic{"loc": {"1-test": gopium.Struct{Name: "test1"}}},
			f: Flat{"1-test": gopium.Struct{Name: "test1"}},
		},
		"single loc multiple items hierarchic collection should return single item flat collection": {
			h: Hierarchic{"loc": {
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			}},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			},
		},
		"multiple locs single item hierarchic collection should return multiple items flat collection": {
			h: Hierarchic{
				"loc-1": {"1-test": gopium.Struct{Name: "test1"}},
				"loc-2": {"2-test": gopium.Struct{Name: "test2"}},
			},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			},
		},
		"multiple locs multiple items hierarchic collection should return multiple items flat collection": {
			h: Hierarchic{
				"loc-1": {
					"1-test": gopium.Struct{Name: "test1"},
					"2-test": gopium.Struct{Name: "test2"},
				},
				"loc-2": {
					"3-test": gopium.Struct{Name: "test3"},
					"4-test": gopium.Struct{Name: "test4"},
				},
			},
			f: Flat{
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
			f := tcase.h.Flat()
			// check
			if !reflect.DeepEqual(f, tcase.f) {
				t.Errorf("actual %v doesn't equal to %v", f, tcase.f)
			}
		})
	}
}
