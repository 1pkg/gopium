package typepkg

import (
	"context"
	"errors"
	"go/token"
	"go/types"
	"reflect"
	"testing"
)

func TestNewMavenGoTypes(t *testing.T) {
	// prepare
	table := map[string]struct {
		compiler string
		arch     string
		caches   []int64
		maven    MavenGoTypes
		err      error
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
	// prepare
	maven, err := NewMavenGoTypes("gc", "amd64", 2, 4, 8, 16, 32)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
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
	// prepare
	maven, err := NewMavenGoTypes("gc", "amd64", 2, 4, 8, 16, 32)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		tp    types.Type
		ctx   context.Context
		name  string
		size  int64
		align int64
	}{
		"int64 type should return correct value": {
			tp:    types.Typ[types.Int64],
			name:  "int64",
			size:  8,
			align: 8,
		},
		"string type should return correct value": {
			tp:    types.Typ[types.String],
			name:  "string",
			size:  16,
			align: 8,
		},
		"string slice type should return correct value": {
			tp:    types.NewSlice(types.Typ[types.String]),
			name:  "[]string",
			size:  24,
			align: 8,
		},
		"float32 arr type should return correct value": {
			tp:    types.NewArray(types.Typ[types.Float32], 8),
			name:  "[8]float32",
			size:  32,
			align: 4,
		},
		"struct type should return correct value": {
			tp: types.NewStruct(
				[]*types.Var{
					types.NewVar(token.Pos(0), nil, "a", types.Typ[types.Int64]),
					types.NewVar(token.Pos(0), nil, "b", types.NewSlice(types.Typ[types.Int64])),
					types.NewVar(token.Pos(0), nil, "c", types.Typ[types.Complex128]),
					types.NewVar(token.Pos(0), nil, "d", types.NewArray(types.Typ[types.Byte], 16)),
				},
				nil,
			),
			name:  "struct{a int64; b []int64; c complex128; d [16]uint8}",
			size:  64,
			align: 8,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			name := maven.Name(tcase.tp)
			size := maven.Size(tcase.tp)
			align := maven.Align(tcase.tp)
			// check
			if name != tcase.name {
				t.Errorf("actual %v doesn't equal to %v", name, tcase.name)
			}
			if size != tcase.size {
				t.Errorf("actual %v doesn't equal to %v", size, tcase.size)
			}
			if align != tcase.align {
				t.Errorf("actual %v doesn't equal to %v", align, tcase.align)
			}
		})
	}

}
