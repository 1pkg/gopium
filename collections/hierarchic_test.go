package collections

import (
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestHierarchicPushCatMixed(t *testing.T) {
	// prepare
	h := NewHierarchic("prefix/")
	h.Push("test-1", "prefix/test-1", gopium.Struct{Name: "test-1"})
	h.Push("test-2", "prefix/test-2", gopium.Struct{Name: "test-2"})
	h.Push("test-3", "prefix/test-2", gopium.Struct{Name: "test-3"})
	h.Push("test-4", "test-2", gopium.Struct{Name: "test-4"})
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
		"prefix/test-1 cat should return expected flat collection": {
			cat: "prefix/test-1",
			f:   Flat{"test-1": gopium.Struct{Name: "test-1"}},
			ok:  true,
		},
		"test-2 cat should return expected flat collection": {
			cat: "test-2",
			f: Flat{
				"test-2": gopium.Struct{Name: "test-2"},
				"test-3": gopium.Struct{Name: "test-3"},
				"test-4": gopium.Struct{Name: "test-4"},
			},
			ok: true,
		},
		"prefix/test-2 cat should return expected flat collection": {
			cat: "prefix/test-2",
			f: Flat{
				"test-2": gopium.Struct{Name: "test-2"},
				"test-3": gopium.Struct{Name: "test-3"},
				"test-4": gopium.Struct{Name: "test-4"},
			},
			ok: true,
		},
		"p/test-2 cat should return empty flat collection": {
			cat: "p/test-2",
			f:   Flat(nil),
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
		h   Hierarchic
		f   Flat
		len int
	}{
		"empty hierarchic collection should return empty flat collection": {
			h:   NewHierarchic(""),
			f:   Flat{},
			len: 0,
		},
		"single cat single item hierarchic collection should return single item flat collection": {
			h:   Hierarchic{cats: map[string]Flat{"loc": {"1-test": gopium.Struct{Name: "test1"}}}},
			f:   Flat{"1-test": gopium.Struct{Name: "test1"}},
			len: 1,
		},
		"single cat multiple items hierarchic collection should return single item flat collection": {
			h: Hierarchic{cats: map[string]Flat{"loc": {
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			}}},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			},
			len: 2,
		},
		"multiple cats single item hierarchic collection should return multiple items flat collection": {
			h: Hierarchic{cats: map[string]Flat{
				"loc-1": {"1-test": gopium.Struct{Name: "test1"}},
				"loc-2": {"2-test": gopium.Struct{Name: "test2"}},
			}},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
			},
			len: 2,
		},
		"multiple cats multiple items hierarchic collection should return multiple items flat collection": {
			h: Hierarchic{cats: map[string]Flat{
				"loc-1": {
					"1-test": gopium.Struct{Name: "test1"},
					"2-test": gopium.Struct{Name: "test2"},
				},
				"loc-2": {
					"3-test": gopium.Struct{Name: "test3"},
					"4-test": gopium.Struct{Name: "test4"},
				},
			}},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
				"3-test": gopium.Struct{Name: "test3"},
				"4-test": gopium.Struct{Name: "test4"},
			},
			len: 4,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			f := tcase.h.Flat()
			len := tcase.h.Len()
			// check
			if !reflect.DeepEqual(f, tcase.f) {
				t.Errorf("actual %v doesn't equal to %v", f, tcase.f)
			}
			if !reflect.DeepEqual(len, tcase.len) {
				t.Errorf("actual %v doesn't equal to %v", len, tcase.len)
			}
		})
	}
}
