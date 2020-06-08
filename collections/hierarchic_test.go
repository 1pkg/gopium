package collections

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"1pkg/gopium/gopium"
)

func TestHierarchicPushCat(t *testing.T) {
	// prepare
	h := NewHierarchic("prefix" + string(os.PathSeparator))
	h.Push("test-1", filepath.Join("prefix", "test-1"), gopium.Struct{Name: "test-1"})
	h.Push("test-2", filepath.Join("prefix", "test-2"))
	h.Push("test-3", filepath.Join("prefix", "test-2"), gopium.Struct{Name: "test-2"}, gopium.Struct{Name: "test-3"})
	h.Push("test-4", "test-2", gopium.Struct{Name: "test-4"})
	h.Push("test-5", "test-3")
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
			cat: filepath.Join("prefix", "test-1"),
			f:   Flat{"test-1": gopium.Struct{Name: "test-1"}},
			ok:  true,
		},
		"test-2 cat should return expected flat collection": {
			cat: "test-2",
			f: Flat{
				"test-3-0": gopium.Struct{Name: "test-2"},
				"test-3-1": gopium.Struct{Name: "test-3"},
				"test-4":   gopium.Struct{Name: "test-4"},
			},
			ok: true,
		},
		"prefix/test-2 cat should return expected flat collection": {
			cat: filepath.Join("prefix", "test-2"),
			f: Flat{
				"test-3-0": gopium.Struct{Name: "test-2"},
				"test-3-1": gopium.Struct{Name: "test-3"},
				"test-4":   gopium.Struct{Name: "test-4"},
			},
			ok: true,
		},
		"p/test-2 cat should return empty flat collection": {
			cat: filepath.Join("p", "test-1"),
			f:   Flat(nil),
		},
		"test-3 cat should return empty flat collection": {
			cat: "test-3",
			f:   make(Flat),
			ok:  true,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			f, ok := h.Catflat(tcase.cat)
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

func TestHierarchicFlatRcatLen(t *testing.T) {
	// prepare
	table := map[string]struct {
		h    Hierarchic
		f    Flat
		rcat string
		len  int
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
		"multiple overlapping cats multiple items hierarchic collection should return multiple items flat collection": {
			h: Hierarchic{cats: map[string]Flat{
				filepath.Join("loc", "test", "123"): {
					"1-test": gopium.Struct{Name: "test1"},
					"2-test": gopium.Struct{Name: "test2"},
				},
				filepath.Join("loc", "test", "abcd"): {
					"3-test": gopium.Struct{Name: "test3"},
					"4-test": gopium.Struct{Name: "test4"},
				},
				filepath.Join("loc", "test", "test", "test", "test"): {
					"5-test": gopium.Struct{Name: "test5"},
				},
				filepath.Join("loc", "test", "abcd", "test", "123"): {
					"6-test": gopium.Struct{Name: "test6"},
				},
			}},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
				"3-test": gopium.Struct{Name: "test3"},
				"4-test": gopium.Struct{Name: "test4"},
				"5-test": gopium.Struct{Name: "test5"},
				"6-test": gopium.Struct{Name: "test6"},
			},
			rcat: filepath.Join("loc", "test"),
			len:  6,
		},
		"multiple non overlapping cats multiple items hierarchic collection should return multiple items flat collection": {
			h: Hierarchic{cats: map[string]Flat{
				filepath.Join("loc1", "test", "123"): {
					"1-test": gopium.Struct{Name: "test1"},
					"2-test": gopium.Struct{Name: "test2"},
				},
				filepath.Join("loc", "test2", "abcd"): {
					"3-test": gopium.Struct{Name: "test3"},
					"4-test": gopium.Struct{Name: "test4"},
				},
				filepath.Join("loc3", "test", "test", "test", "test"): {
					"5-test": gopium.Struct{Name: "test5"},
				},
				filepath.Join("loc", "test", "abcd", "test", "123"): {
					"6-test": gopium.Struct{Name: "test6"},
				},
			}},
			f: Flat{
				"1-test": gopium.Struct{Name: "test1"},
				"2-test": gopium.Struct{Name: "test2"},
				"3-test": gopium.Struct{Name: "test3"},
				"4-test": gopium.Struct{Name: "test4"},
				"5-test": gopium.Struct{Name: "test5"},
				"6-test": gopium.Struct{Name: "test6"},
			},
			len: 6,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			f := tcase.h.Flat()
			hlen := tcase.h.Len()
			rcat := tcase.h.Rcat()
			// check
			if !reflect.DeepEqual(f, tcase.f) {
				t.Errorf("actual %v doesn't equal to %v", f, tcase.f)
			}
			if !reflect.DeepEqual(hlen, tcase.len) {
				t.Errorf("actual %v doesn't equal to %v", hlen, tcase.len)
			}
			if !reflect.DeepEqual(rcat, tcase.rcat) {
				t.Errorf("actual %v doesn't equal to %v", rcat, tcase.rcat)
			}
		})
	}
}
