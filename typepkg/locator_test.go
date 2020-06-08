package typepkg

import (
	"go/token"
	"reflect"
	"testing"

	"1pkg/gopium/gopium"
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
			id:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08:1",
			loc: "test",
		},
		"valid token pos 11 should be located in expected file": {
			pos: token.Pos(11),
			id:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08:3",
			loc: "test",
		},
		"valid token pos 21 should be located in expected file": {
			pos: token.Pos(21),
			id:  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08:3",
			loc: "test",
		},
		"valid token pos 22 should be located in expected file": {
			pos: token.Pos(22),
			id:  "7395c92698585f5dff68edce8314a30f534e98707af84f11a725cd62b538b139:1",
			loc: "loc-test-id",
		},
		"valid token pos 50 should be located in expected file": {
			pos: token.Pos(50),
			id:  "7395c92698585f5dff68edce8314a30f534e98707af84f11a725cd62b538b139:2",
			loc: "loc-test-id",
		},
		"valid token pos 52 should be located in expected file": {
			pos: token.Pos(52),
			id:  "7395c92698585f5dff68edce8314a30f534e98707af84f11a725cd62b538b139:2",
			loc: "loc-test-id",
		},
		"valid token pos 53 should be located in expected file": {
			pos: token.Pos(53),
			id:  "e65de074e1a2fc5a98d431c9d737f851e27625eca54b73816d4682805938e454:1",
			loc: "id-test-loc",
		},
		"valid token pos 99 should be located in expected file": {
			pos: token.Pos(99),
			id:  "e65de074e1a2fc5a98d431c9d737f851e27625eca54b73816d4682805938e454:3",
			loc: "id-test-loc",
		},
		"valid token pos 100 should be located in expected file": {
			pos: token.Pos(100),
			id:  "e65de074e1a2fc5a98d431c9d737f851e27625eca54b73816d4682805938e454:3",
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
