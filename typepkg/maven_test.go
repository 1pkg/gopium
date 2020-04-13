package typepkg

import (
	"errors"
	"go/token"
	"go/types"
	"reflect"
	"testing"
)

func TestNewMavenGoTypes(t *testing.T) {
	// prepare
	table := map[string]struct {
		compiler, arch string
		caches         []int64
		maven          MavenGoTypes
		err            error
	}{
		"non existed compiler should spawn error": {
			compiler: "test",
			arch:     "amd64",
			maven:    MavenGoTypes{},
			err:      errors.New(`unsuported compiler "test" arch "amd64" combination`),
		},
		"non existed arch should spawn error": {
			compiler: "gc",
			arch:     "test",
			maven:    MavenGoTypes{},
			err:      errors.New(`unsuported compiler "gc" arch "test" combination`),
		},
		"existed compiler and arch should create actual maven": {
			compiler: "gc",
			arch:     "amd64",
			caches:   []int64{2, 4, 8, 16, 32},
			maven: MavenGoTypes{
				sizes: types.SizesFor("gc", "amd64"),
				caches: map[uint]int64{
					1: 2,
					2: 4,
					3: 8,
					4: 16,
					5: 32,
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			maven, err := NewMavenGoTypes(tcase.compiler, tcase.arch, tcase.caches...)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			if !reflect.DeepEqual(maven, tcase.maven) {
				t.Errorf("actual %v doesn't equal to expected %v", maven, tcase.maven)
			}
		})
	}
}

func TestMavenGoTypesSysMixed(t *testing.T) {
	maven, err := NewMavenGoTypes("gc", "amd64", 2, 4, 8, 16, 32)
	if err != nil {
		t.Errorf("actual %v doesn't equal to %v", err, nil)
	}
	// check sys word
	if maven.SysWord() != 8 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysWord(), 8)
	}
	// check sys align
	if maven.SysAlign() != 8 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysAlign(), 8)
	}
	// check sys cache l1
	if maven.SysCache(1) != 2 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysCache(1), 2)
	}
	// check sys cache l2
	if maven.SysCache(2) != 4 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysCache(2), 4)
	}
	// check sys cache l3
	if maven.SysCache(3) != 8 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysCache(3), 8)
	}
	// check sys cache l10
	if maven.SysCache(10) != 64 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysCache(10), 64)
	}
	// check sys cache l0
	if maven.SysCache(0) != 64 {
		t.Errorf("actual %v doesn't equal to %v", maven.SysCache(0), 64)
	}
}

func TestMavenGoTypesTypeMixed(t *testing.T) {
	maven, err := NewMavenGoTypes("gc", "amd64", 2, 4, 8, 16, 32)
	if err != nil {
		t.Errorf("actual %v doesn't equal to %v", err, nil)
	}
	var tp types.Type
	// check int64
	tp = types.Typ[types.Int64]
	name := maven.Name(tp)
	if name != "int64" {
		t.Errorf("actual %v doesn't equal to %v", name, "int64")
	}
	size := maven.Size(tp)
	if size != 8 {
		t.Errorf("actual %v doesn't equal to %v", size, 8)
	}
	align := maven.Align(tp)
	if align != 8 {
		t.Errorf("actual %v doesn't equal to %v", align, 8)
	}
	// check string
	tp = types.Typ[types.String]
	name = maven.Name(tp)
	if name != "string" {
		t.Errorf("actual %v doesn't equal to %v", name, "string")
	}
	size = maven.Size(tp)
	if size != 16 {
		t.Errorf("actual %v doesn't equal to %v", size, 16)
	}
	align = maven.Align(tp)
	if align != 8 {
		t.Errorf("actual %v doesn't equal to %v", align, 8)
	}
	// check string slice
	tp = types.Typ[types.String]
	tp = types.NewSlice(tp)
	name = maven.Name(tp)
	if name != "[]string" {
		t.Errorf("actual %v doesn't equal to %v", name, "[]string")
	}
	size = maven.Size(tp)
	if size != 24 {
		t.Errorf("actual %v doesn't equal to %v", size, 24)
	}
	align = maven.Align(tp)
	if align != 8 {
		t.Errorf("actual %v doesn't equal to %v", align, 8)
	}
	// check float32 arr
	tp = types.Typ[types.Float32]
	tp = types.NewArray(tp, 8)
	name = maven.Name(tp)
	if name != "[8]float32" {
		t.Errorf("actual %v doesn't equal to %v", name, "[8]float32")
	}
	size = maven.Size(tp)
	if size != 32 {
		t.Errorf("actual %v doesn't equal to %v", size, 32)
	}
	align = maven.Align(tp)
	if align != 4 {
		t.Errorf("actual %v doesn't equal to %v", align, 4)
	}
	// check custom struct
	tp = types.NewStruct(
		[]*types.Var{
			types.NewVar(token.Pos(0), nil, "a", types.Typ[types.Int64]),
			types.NewVar(token.Pos(0), nil, "b", types.NewSlice(types.Typ[types.Int64])),
			types.NewVar(token.Pos(0), nil, "c", types.Typ[types.Complex128]),
			types.NewVar(token.Pos(0), nil, "d", types.NewArray(types.Typ[types.Byte], 16)),
		},
		nil,
	)
	name = maven.Name(tp)
	if name != "struct{a int64; b []int64; c complex128; d [16]uint8}" {
		t.Errorf("actual %v doesn't equal to %v", name, "struct{a int64; b []int64; c complex128; d [16]uint8}")
	}
	size = maven.Size(tp)
	if size != 64 {
		t.Errorf("actual %v doesn't equal to %v", size, 64)
	}
	align = maven.Align(tp)
	if align != 8 {
		t.Errorf("actual %v doesn't equal to %v", align, 8)
	}
}
