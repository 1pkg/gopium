package walkers

import (
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests/mocks"
)

func TestMavenHas(t *testing.T) {
	// prepare
	m := maven{
		loc: mocks.Locator{
			Poses: map[token.Pos]mocks.Pos{
				token.Pos(0): {
					ID:  "1",
					Loc: "loc1",
				},
				token.Pos(10): {
					ID:  "10",
					Loc: "loc10",
				},
			},
		},
	}
	table := map[string]struct {
		tn  *types.TypeName
		id  string
		loc string
	}{
		"type name with valid pos should return expected id and loc": {
			tn:  types.NewTypeName(token.Pos(0), nil, "test", types.Typ[types.String]),
			id:  "1",
			loc: "loc1",
		},
		"other type name with valid pos should return expected id and loc": {
			tn:  types.NewTypeName(token.Pos(10), nil, "test", types.Typ[types.String]),
			id:  "10",
			loc: "loc10",
		},
		"type name with invalid pos should provide empty id and loc": {
			tn:  types.NewTypeName(token.Pos(100), nil, "test", types.Typ[types.String]),
			id:  "",
			loc: "",
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			id1, loc1, ok1 := m.has(tcase.tn)
			id2, loc2, ok2 := m.has(tcase.tn)
			// check
			if !reflect.DeepEqual(ok1, false) {
				t.Errorf("actual %v doesn't equal to expected %v", ok1, false)
			}
			if !reflect.DeepEqual(id1, tcase.id) {
				t.Errorf("actual %v doesn't equal to expected %v", id1, tcase.id)
			}
			if !reflect.DeepEqual(loc1, tcase.loc) {
				t.Errorf("actual %v doesn't equal to expected %v", loc1, tcase.loc)
			}
			if !reflect.DeepEqual(ok2, true) {
				t.Errorf("actual %v doesn't equal to expected %v", ok2, true)
			}
			if !reflect.DeepEqual(id2, tcase.id) {
				t.Errorf("actual %v doesn't equal to expected %v", id2, tcase.id)
			}
			if !reflect.DeepEqual(loc2, tcase.loc) {
				t.Errorf("actual %v doesn't equal to expected %v", loc2, tcase.loc)
			}
		})
	}
}

func TestMavenEnum(t *testing.T) {
	// prepare
	ref := collections.NewReference(true)
	ref.Alloc("test")
	ref.Set("test", ptrsizealign{ptr: 32, size: 32, align: 32})
	m := maven{
		exp: mocks.Maven{
			Types: map[string]mocks.Type{
				"string": {
					Name:  "string",
					Size:  16,
					Align: 8,
					Ptr:   16,
				},
				"test": {
					Name:  "test",
					Size:  24,
					Align: 20,
					Ptr:   16,
				},
				"[10]test": {
					Name:  "test",
					Size:  240,
					Align: 20,
					Ptr:   240,
				},
			},
		},
		loc: mocks.Locator{
			Poses: map[token.Pos]mocks.Pos{
				token.Pos(0): {
					ID:  "test",
					Loc: "loc",
				},
			},
		},
		ref: ref,
	}
	sti := types.NewStruct(
		[]*types.Var{
			types.NewVar(token.Pos(0), nil, "a", types.Typ[types.String]),
			types.NewVar(token.Pos(0), nil, "b", types.Typ[types.String]),
			types.NewVar(token.Pos(0), nil, "c", types.Typ[types.String]),
		},
		nil,
	)
	tp := types.NewNamed(types.NewTypeName(token.Pos(0), nil, "test", sti), sti, nil)
	table := map[string]struct {
		name string
		tst  *types.Struct
		st   gopium.Struct
	}{
		"custom type should return expected struct": {
			name: "test-st",
			tst:  sti,
			st: gopium.Struct{
				Name: "test-st",
				Fields: []gopium.Field{
					{
						Name:  "a",
						Type:  "string",
						Size:  16,
						Align: 8,
						Ptr:   16,
					},
					{
						Name:  "b",
						Type:  "string",
						Size:  16,
						Align: 8,
						Ptr:   16,
					},
					{
						Name:  "c",
						Type:  "string",
						Size:  16,
						Align: 8,
						Ptr:   16,
					},
				},
			},
		},
		"custom type with backref should return expected struct": {
			name: "test-st",
			tst:  types.NewStruct([]*types.Var{types.NewVar(token.Pos(0), nil, "v", tp)}, nil),
			st: gopium.Struct{
				Name: "test-st",
				Fields: []gopium.Field{
					{
						Name:  "v",
						Type:  "test",
						Size:  32,
						Align: 32,
						Ptr:   32,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			st := m.enum(tcase.name, tcase.tst)
			// check
			if !reflect.DeepEqual(st, tcase.st) {
				t.Errorf("actual %v doesn't equal to expected %v", st, tcase.st)
			}
		})
	}
}

func TestMavenRefsa(t *testing.T) {
	// prepare
	ref := collections.NewReference(true)
	ref.Alloc("test")
	ref.Set("test", ptrsizealign{ptr: 32, size: 32, align: 32})
	m := maven{
		exp: mocks.Maven{
			Types: map[string]mocks.Type{
				"string": {
					Name:  "string",
					Size:  16,
					Align: 8,
					Ptr:   16,
				},
				"test": {
					Name:  "test",
					Size:  24,
					Align: 20,
					Ptr:   16,
				},
				"[10]test": {
					Name:  "test",
					Size:  240,
					Align: 20,
					Ptr:   240,
				},
			},
		},
		loc: mocks.Locator{
			Poses: map[token.Pos]mocks.Pos{
				token.Pos(0): {
					ID:  "test",
					Loc: "loc",
				},
			},
		},
	}
	st := types.NewStruct([]*types.Var{types.NewVar(token.Pos(0), nil, "a", types.Typ[types.Int64])}, nil)
	tp := types.NewNamed(types.NewTypeName(token.Pos(0), nil, "test", st), st, nil)
	table := map[string]struct {
		t   types.Type
		sa  ptrsizealign
		ref *collections.Reference
	}{
		"primitive type should return expected size and align without backref": {
			t:  types.Typ[types.String],
			sa: ptrsizealign{size: 16, align: 8, ptr: 16},
		},
		"custom type should return expected size and align without backref": {
			t:  tp,
			sa: ptrsizealign{size: 24, align: 20, ptr: 16},
		},
		"custom arr type should return expected size and align without backref": {
			t:  types.NewArray(tp, 10),
			sa: ptrsizealign{size: 240, align: 20, ptr: 240},
		},
		"primitive type should return expected size and align with backref": {
			t:   types.Typ[types.String],
			sa:  ptrsizealign{size: 16, align: 8, ptr: 16},
			ref: ref,
		},
		"custom type should return expected size and align with backref": {
			t:   tp,
			sa:  ptrsizealign{size: 32, align: 32, ptr: 32},
			ref: ref,
		},
		"custom arr type should return expected size and align with backref": {
			t:   types.NewArray(tp, 10),
			sa:  ptrsizealign{size: 320, align: 32, ptr: 320},
			ref: ref,
		},
		"custom non struct type should return expected size and align with backref": {
			t: types.NewNamed(
				types.NewTypeName(token.Pos(0), nil, "test", types.Typ[types.Int64]),
				types.Typ[types.Int64],
				nil,
			),
			sa:  ptrsizealign{size: 24, align: 20, ptr: 16},
			ref: ref,
		},
		"custom empty arr type should return expected size and align with backref": {
			t:   types.NewArray(tp, 0),
			sa:  ptrsizealign{},
			ref: ref,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			ml := maven{exp: m.exp, loc: m.loc}
			ml.ref = tcase.ref
			// exec
			sa := ml.refpsa(tcase.t)
			// check
			if !reflect.DeepEqual(sa, tcase.sa) {
				t.Errorf("actual %v doesn't equal to expected %v", sa, tcase.sa)
			}
		})
	}
}

func TestMavenRefst(t *testing.T) {
	// prepare
	m := maven{ref: collections.NewReference(true)}
	f1, f2, f3 := m.refst("test-1"), m.refst("test-2"), m.refst("test-3")
	f2(gopium.Struct{
		Fields: []gopium.Field{
			{Size: 4, Align: 4, Ptr: 4},
			{Size: 6, Align: 6, Ptr: 6},
			{Size: 2, Align: 2, Ptr: 2},
		},
	})
	go func() {
		f1(gopium.Struct{
			Fields: []gopium.Field{
				{Size: 8, Align: 8, Ptr: 8},
			},
		})
		go func() {
			f3(gopium.Struct{
				Fields: []gopium.Field{
					{Size: 4, Align: 4, Ptr: 4},
					{Size: 6, Align: 4, Ptr: 1},
				},
			})
		}()
	}()
	table := map[string]struct {
		key string
		sa  ptrsizealign
	}{
		"test-1 key should return expected result": {
			key: "test-1",
			sa:  ptrsizealign{ptr: 8, size: 8, align: 8},
		},
		"test-2 key should return expected result": {
			key: "test-2",
			sa:  ptrsizealign{ptr: 14, size: 18, align: 6},
		},
		"test-3 key should return expected result": {
			key: "test-3",
			sa:  ptrsizealign{ptr: 5, size: 12, align: 4},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			sa := m.ref.Get(tcase.key)
			// check
			if !reflect.DeepEqual(sa, tcase.sa) {
				t.Errorf("actual %v doesn't equal to expected %v", sa, tcase.sa)
			}
		})
	}
}
