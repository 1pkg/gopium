package typepkg

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/gopium"
)

func TestNewMavenGoTypes(t *testing.T) {
	// prepare
	table := map[string]struct {
		compiler string
		arch     string
		caches   []int64
		maven    gopium.Maven
		err      error
	}{
		"invalid compiler should return error": {
			compiler: "test",
			arch:     "amd64",
			maven:    MavenGoTypes{},
			err:      errors.New(`unsuported compiler "test" arch "amd64" combination`),
		},
		"invalid arch should return error": {
			compiler: "gc",
			arch:     "test",
			maven:    MavenGoTypes{},
			err:      errors.New(`unsuported compiler "gc" arch "test" combination`),
		},
		"valid compiler and arch pair should return expected maven": {
			compiler: "gc",
			arch:     "amd64",
			caches:   []int64{2, 4, 8, 16, 32},
			maven: MavenGoTypes{
				sizes: stdsizes{types.SizesFor("gc", "amd64").(*types.StdSizes)},
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
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			if !reflect.DeepEqual(maven, tcase.maven) {
				t.Errorf("actual %v doesn't equal to expected %v", maven, tcase.maven)
			}
		})
	}
}

func TestMavenGoTypesCurator(t *testing.T) {
	// prepare
	maven, err := NewMavenGoTypes("gc", "amd64", 2, 4, 8, 16, 32)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		maven  gopium.Maven
		word   int64
		align  int64
		caches []int64
	}{
		"gc/amd64 maven should return expected results": {
			maven:  maven,
			word:   8,
			align:  8,
			caches: []int64{64, 2, 4, 8, 16, 32, 64, 64, 64},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			word := maven.SysWord()
			align := maven.SysAlign()
			caches := make([]int64, len(tcase.caches))
			for i := range caches {
				caches[i] = maven.SysCache(uint(i))
			}
			// check
			if !reflect.DeepEqual(word, tcase.word) {
				t.Errorf("actual %v doesn't equal to %v", word, tcase.word)
			}
			if !reflect.DeepEqual(align, tcase.align) {
				t.Errorf("actual %v doesn't equal to %v", align, tcase.align)
			}
			if !reflect.DeepEqual(caches, tcase.caches) {
				t.Errorf("actual %v doesn't equal to %v", caches, tcase.caches)
			}
		})
	}
}

func TestMavenGoTypesExposer(t *testing.T) {
	// prepare
	maven, err := NewMavenGoTypes("gc", "amd64", 2, 4, 8, 16, 32)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		tp    types.Type
		name  string
		size  int64
		align int64
		ptr   int64
	}{
		"int64 type should return expected resultss": {
			tp:    types.Typ[types.Int64],
			name:  "int64",
			size:  8,
			align: 8,
			ptr:   0,
		},
		"string type should return expected results": {
			tp:    types.Typ[types.String],
			name:  "string",
			size:  16,
			align: 8,
			ptr:   8,
		},
		"string slice type should return expected results": {
			tp:    types.NewSlice(types.Typ[types.String]),
			name:  "[]string",
			size:  24,
			align: 8,
			ptr:   8,
		},
		"float32 arr type should return expected results": {
			tp:    types.NewArray(types.Typ[types.Float32], 8),
			name:  "[8]float32",
			size:  32,
			align: 4,
			ptr:   0,
		},
		"struct type should return expected results": {
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
			ptr:   16,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			name := maven.Name(tcase.tp)
			size := maven.Size(tcase.tp)
			align := maven.Align(tcase.tp)
			ptr := maven.Ptr(tcase.tp)
			// check
			if !reflect.DeepEqual(name, tcase.name) {
				t.Errorf("actual %v doesn't equal to %v", name, tcase.name)
			}
			if !reflect.DeepEqual(size, tcase.size) {
				t.Errorf("actual %v doesn't equal to %v", size, tcase.size)
			}
			if !reflect.DeepEqual(align, tcase.align) {
				t.Errorf("actual %v doesn't equal to %v", align, tcase.align)
			}
			if !reflect.DeepEqual(ptr, tcase.ptr) {
				t.Errorf("actual %v doesn't equal to %v", ptr, tcase.ptr)
			}
		})
	}

}
