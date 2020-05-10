package typepkg

import (
	"go/token"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestNewLocatorRoot(t *testing.T) {
	// prepare
	fset := token.NewFileSet()
	fset.AddFile("test", 1, 10)
	table := map[string]struct {
		fset *token.FileSet
		root *token.FileSet
		loc  *Locator
	}{
		"nil fset should return default locator": {
			root: token.NewFileSet(),
			loc: &Locator{
				root:  token.NewFileSet(),
				extra: make(map[string]*token.FileSet),
			},
		},
		"non nil fset should return custom locator": {
			fset: fset,
			root: fset,
			loc: &Locator{
				root:  fset,
				extra: make(map[string]*token.FileSet),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			loc := NewLocator(tcase.fset)
			root := loc.Root()
			// check
			if !reflect.DeepEqual(loc, tcase.loc) {
				t.Errorf("actual %v doesn't equal to expected %v", loc, tcase.loc)
			}
			if !reflect.DeepEqual(root, tcase.root) {
				t.Errorf("actual %v doesn't equal to expected %v", root, tcase.root)
			}
		})
	}
}

func TestLocatorIDLoc(t *testing.T) {
	// prepare
	fset := token.NewFileSet()
	f1 := fset.AddFile("test", 1, 20)
	f1.AddLine(5)
	f1.AddLine(10)
	f2 := fset.AddFile("loc-test-id", 22, 30)
	f2.AddLine(20)
	f3 := fset.AddFile("id-test-loc", 53, 47)
	f3.AddLine(1)
	f3.AddLine(10)
	locator := NewLocator(fset)
	table := map[string]struct {
		pos token.Pos
		id  string
		loc string
	}{
		"valid token pos 1 should be located in expected file": {
			pos: token.Pos(1),
			id:  "1-cad4a5be62ba01bfe7a07a8ff9ab1ed0d726c3cd82bfb3053f92fc21b3088ca3",
			loc: "test",
		},
		"valid token pos 11 should be located in expected file": {
			pos: token.Pos(11),
			id:  "3-da1d0a859e4d55d60b29d1a8b8ce379a9c24b7e1db83868708329c64193470bb",
			loc: "test",
		},
		"valid token pos 21 should be located in expected file": {
			pos: token.Pos(21),
			id:  "3-da1d0a859e4d55d60b29d1a8b8ce379a9c24b7e1db83868708329c64193470bb",
			loc: "test",
		},
		"valid token pos 22 should be located in expected file": {
			pos: token.Pos(22),
			id:  "1-a79ce52b40bfe7dfc16f512d45d9d382cefb70603f50adb5abcf5f73f4b4fefe",
			loc: "loc-test-id",
		},
		"valid token pos 50 should be located in expected file": {
			pos: token.Pos(50),
			id:  "2-4a7e3c3497f71fdb5c1c2389cd6d2e6afb93706c72448b8308e643b4ab56a791",
			loc: "loc-test-id",
		},
		"valid token pos 52 should be located in expected file": {
			pos: token.Pos(52),
			id:  "2-4a7e3c3497f71fdb5c1c2389cd6d2e6afb93706c72448b8308e643b4ab56a791",
			loc: "loc-test-id",
		},
		"valid token pos 53 should be located in expected file": {
			pos: token.Pos(53),
			id:  "1-80b7343d7bde2f986326d4d4b6c638b24f22f3a46b7e1f1eac80488e90f91398",
			loc: "id-test-loc",
		},
		"valid token pos 99 should be located in expected file": {
			pos: token.Pos(99),
			id:  "3-f898adf4c5d8f97ed4f7841d2afacb8225690dd19a538b08071f50d20d44f79c",
			loc: "id-test-loc",
		},
		"valid token pos 100 should be located in expected file": {
			pos: token.Pos(100),
			id:  "3-f898adf4c5d8f97ed4f7841d2afacb8225690dd19a538b08071f50d20d44f79c",
			loc: "id-test-loc",
		},
		"invalid token pos 1000 should return default results": {
			pos: token.Pos(1000),
			id:  "",
			loc: "",
		},
		"invalid token pos -1 should return default results": {
			pos: token.Pos(-1),
			id:  "",
			loc: "",
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			id, loc := locator.ID(tcase.pos), locator.Loc(tcase.pos)
			// check
			if !reflect.DeepEqual(id, tcase.id) {
				t.Errorf("actual %v doesn't equal to expected %v", id, tcase.id)
			}
			if !reflect.DeepEqual(loc, tcase.loc) {
				t.Errorf("actual %v doesn't equal to expected %v", loc, tcase.loc)
			}
		})
	}
}

func TestLocatorFset(t *testing.T) {
	// prepare
	fset := token.NewFileSet()
	fset.AddFile("test", 1, 10)
	locator := NewLocator(fset)
	tfset := token.NewFileSet()
	tfset, ok := locator.Fset("test", tfset)
	if !reflect.DeepEqual(ok, true) {
		t.Fatalf("actual %v doesn't equal to %v", ok, true)
	}
	if reflect.DeepEqual(tfset, nil) {
		t.Fatalf("actual %v doesn't equal to not %v", tfset, nil)
	}
	table := map[string]struct {
		l      string
		loc    gopium.Locator
		fset   *token.FileSet
		okloc  bool
		okfset bool
	}{
		"invalid loc should return default results": {
			l:    "loc",
			loc:  NewLocator(fset),
			fset: fset,
		},
		"valid loc should return expected results": {
			l:      "test",
			loc:    NewLocator(tfset),
			fset:   tfset,
			okloc:  true,
			okfset: true,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			loc, okloc := locator.Locator(tcase.l)
			fset, okfset := locator.Fset(tcase.l, nil)
			// check
			if !reflect.DeepEqual(okloc, tcase.okloc) {
				t.Errorf("actual %v doesn't equal to expected %v", okloc, tcase.okloc)
			}
			if !reflect.DeepEqual(loc, tcase.loc) {
				t.Errorf("actual %v doesn't equal to expected %v", loc, tcase.loc)
			}
			if !reflect.DeepEqual(okfset, tcase.okfset) {
				t.Errorf("actual %v doesn't equal to expected %v", okfset, tcase.okfset)
			}
			if !reflect.DeepEqual(fset, tcase.fset) {
				t.Errorf("actual %v doesn't equal to expected %v", fset, tcase.fset)
			}
		})
	}
}
