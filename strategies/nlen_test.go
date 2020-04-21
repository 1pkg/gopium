package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestNLen(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		nlen nlen
		ctx  context.Context
		o    gopium.Struct
		r    gopium.Struct
		err  error
	}{
		"empty struct should be applied to empty struct": {
			nlen: nlenasc,
			ctx:  context.Background(),
		},
		"non empty struct should be applied to itself": {
			nlen: nlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
		},
		"non empty struct should be applied to itself on canceled context": {
			nlen: nlendesc,
			ctx:  cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
					},
				},
			},
			err: cctx.Err(),
		},
		"asc name len struct should be applied to itself": {
			nlen: nlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "atest1-t",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "ztest3-ttt",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "atest1-t",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "ztest3-ttt",
					},
				},
			},
		},
		"desc name len struct should be applied to sorted struct": {
			nlen: nlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "ztest1-ttt",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "atest3-t",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "atest3-t",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "ztest1-ttt",
					},
				},
			},
		},
		"mixed name len struct should be applied to sorted struct asc": {
			nlen: nlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test-1",
					},
					{
						Name: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name: "test1000000",
					},
					{
						Name:  "zzz",
						Type:  "int33",
						Size:  10,
						Align: 10,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "zzz",
						Type:  "int33",
						Size:  10,
						Align: 10,
					},
					{
						Name: "test-1",
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name: "test1000000",
					},
				},
			},
		},
		"asc name len struct should be applied to sorted struct": {
			nlen: nlendesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "atest1-t",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "ztest3-ttt",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "ztest3-ttt",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "atest1-t",
					},
				},
			},
		},
		"desc name len struct should be applied to itself": {
			nlen: nlendesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "ztest1-ttt",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "atest3-t",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "ztest1-ttt",
					},
					{
						Name: "rtest2-tt",
					},
					{
						Name: "atest3-t",
					},
				},
			},
		},
		"mixed name len struct should be applied to sorted struct desc": {
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test-1",
					},
					{
						Name: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name: "test1000000",
					},
					{
						Name:  "zzz",
						Type:  "int33",
						Size:  10,
						Align: 10,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1000000",
					},
					{
						Name: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name: "test-1",
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name:  "zzz",
						Type:  "int33",
						Size:  10,
						Align: 10,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.nlen.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
