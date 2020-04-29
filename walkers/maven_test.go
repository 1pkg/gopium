package walkers

import (
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/tests/mocks"
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
		"type name with valid pos should provide correct id and loc": {
			tn:  types.NewTypeName(token.Pos(0), nil, "test", types.Typ[types.String]),
			id:  "1",
			loc: "loc1",
		},
		"different type name with valid pos should provide correct id and loc": {
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
			id, loc, ok := m.has(tcase.tn)
			id_, loc_, ok_ := m.has(tcase.tn)
			// check
			if ok || id != tcase.id {
				t.Errorf("actual %v doesn't equal to expected %v", id, tcase.id)
			}
			if ok || loc != tcase.loc {
				t.Errorf("actual %v doesn't equal to expected %v", loc, tcase.loc)
			}
			if !ok_ || id_ != tcase.id {
				t.Errorf("actual %v doesn't equal to expected %v", id_, tcase.id)
			}
			if !ok_ || loc_ != tcase.loc {
				t.Errorf("actual %v doesn't equal to expected %v", loc_, tcase.loc)
			}
		})
	}
}

func TestMavenEnum(t *testing.T) {
	// prepare
	ref := collections.NewReference(false)
	ref.Alloc("test")
	ref.Set("test", sizealign{size: 32, align: 32})
	m := maven{
		exp: mocks.Maven{
			Types: map[string]mocks.Type{
				"string": mocks.Type{
					Name:  "string",
					Size:  16,
					Align: 8,
				},
				"test": mocks.Type{
					Name:  "test",
					Size:  24,
					Align: 20,
				},
				"[10]test": mocks.Type{
					Name:  "test",
					Size:  240,
					Align: 20,
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
		"custom type should return correct enum struct": {
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
					},
					{
						Name:  "b",
						Type:  "string",
						Size:  16,
						Align: 8,
					},
					{
						Name:  "c",
						Type:  "string",
						Size:  16,
						Align: 8,
					},
				},
			},
		},
		"custom type from backref should return correct enum struct": {
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
	ref := collections.NewReference(false)
	ref.Alloc("test")
	ref.Set("test", sizealign{size: 32, align: 32})
	m := maven{
		exp: mocks.Maven{
			Types: map[string]mocks.Type{
				"string": mocks.Type{
					Name:  "string",
					Size:  16,
					Align: 8,
				},
				"test": mocks.Type{
					Name:  "test",
					Size:  24,
					Align: 20,
				},
				"[10]test": mocks.Type{
					Name:  "test",
					Size:  240,
					Align: 20,
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
		sa  sizealign
		ref *collections.Reference
	}{
		"primitive type should return correct size and align without backref": {
			t:  types.Typ[types.String],
			sa: sizealign{size: 16, align: 8},
		},
		"custom type should return correct size and align without backref": {
			t:  tp,
			sa: sizealign{size: 24, align: 20},
		},
		"custom arr type should return correct size and align without backref": {
			t:  types.NewArray(tp, 10),
			sa: sizealign{size: 240, align: 20},
		},
		"primitive type should return correct size and align with actual backref": {
			t:   types.Typ[types.String],
			sa:  sizealign{size: 16, align: 8},
			ref: ref,
		},
		"custom type should return correct size and align with actual backref": {
			t:   tp,
			sa:  sizealign{size: 32, align: 32},
			ref: ref,
		},
		"custom arr type should return correct size and align with actual backref": {
			t:   types.NewArray(tp, 10),
			sa:  sizealign{size: 320, align: 32},
			ref: ref,
		},
		"custom non struct type should return correct size and align with actual backref": {
			t: types.NewNamed(
				types.NewTypeName(token.Pos(0), nil, "test", types.Typ[types.Int64]),
				types.Typ[types.Int64],
				nil,
			),
			sa:  sizealign{size: 24, align: 20},
			ref: ref,
		},
		"custom empty arr type should return correct size and align with actual backref": {
			t:   types.NewArray(tp, 0),
			sa:  sizealign{},
			ref: ref,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			ml := maven{exp: m.exp, loc: m.loc}
			ml.ref = tcase.ref
			sa := ml.refsa(tcase.t)
			// check
			if !reflect.DeepEqual(sa, tcase.sa) {
				t.Errorf("actual %v doesn't equal to expected %v", sa, tcase.sa)
			}
		})
	}
}

func TestMavenRefst(t *testing.T) {
	// prepare
	m := maven{ref: collections.NewReference(false)}
	// test with real ref
	f1, f2 := m.refst("test-1"), m.refst("test-2")
	f2(gopium.Struct{
		Fields: []gopium.Field{
			{Size: 4, Align: 4},
			{Size: 6, Align: 6},
			{Size: 2, Align: 2},
		},
	})
	// set f1 in goroutine
	go func() {
		f1(gopium.Struct{
			Fields: []gopium.Field{
				{Size: 8, Align: 8},
			},
		})
	}()
	sa1, sa2 := m.ref.Get("test-1"), m.ref.Get("test-2")
	if !reflect.DeepEqual(sa1, sizealign{size: 8, align: 8}) {
		t.Errorf("actual %v doesn't equal to expected %v", sa1, sizealign{size: 8, align: 8})
	}
	if !reflect.DeepEqual(sa2, sizealign{size: 12, align: 6}) {
		t.Errorf("actual %v doesn't equal to expected %v", sa2, sizealign{size: 12, align: 6})
	}
	f2(gopium.Struct{
		Fields: []gopium.Field{
			{Size: 4, Align: 4},
			{Size: 6, Align: 4},
		},
	})
	sa1, sa2 = m.ref.Get("test-1"), m.ref.Get("test-2")
	if !reflect.DeepEqual(sa1, sizealign{size: 8, align: 8}) {
		t.Errorf("actual %v doesn't equal to expected %v", sa1, sizealign{size: 8, align: 8})
	}
	if !reflect.DeepEqual(sa2, sizealign{size: 10, align: 4}) {
		t.Errorf("actual %v doesn't equal to expected %v", sa2, sizealign{size: 10, align: 4})
	}
}
