package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestTLen(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		tlen tlen
		ctx  context.Context
		o    gopium.Struct
		r    gopium.Struct
		err  error
	}{
		"empty struct should be applied to empty struct": {
			tlen: tlenasc,
			ctx:  context.Background(),
		},
		"non empty struct should be applied to itself": {
			tlen: tlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test",
					},
				},
			},
		},
		"non empty struct should be applied to itself on canceled context": {
			tlen: tlendesc,
			ctx:  cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test",
					},
				},
			},
			err: cctx.Err(),
		},
		"asc type len struct should be applied to itself": {
			tlen: tlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "atest1-t",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "ztest3-ttt",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "atest1-t",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "ztest3-ttt",
					},
				},
			},
		},
		"desc type len struct should be applied to sorted struct": {
			tlen: tlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "ztest1-ttt",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "atest3-t",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "atest3-t",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "ztest1-ttt",
					},
				},
			},
		},
		"mixed type len struct should be applied to sorted struct asc": {
			tlen: tlenasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test-1",
					},
					{
						Name: "test",
						Type: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name: "test",
						Type: "test1000000",
					},
					{
						Name:  "test",
						Type:  "zzz",
						Size:  10,
						Align: 10,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Type:  "zzz",
						Size:  10,
						Align: 10,
					},
					{
						Name: "test",
						Type: "test-1",
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name: "test",
						Type: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name: "test",
						Type: "test1000000",
					},
				},
			},
		},
		"asc type len struct should be applied to sorted struct": {
			tlen: tlendesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "atest1-t",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "ztest3-ttt",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "ztest3-ttt",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "atest1-t",
					},
				},
			},
		},
		"desc type len struct should be applied to itself": {
			tlen: tlendesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "ztest1-ttt",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "atest3-t",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "ztest1-ttt",
					},
					{
						Name: "test",
						Type: "rtest2-tt",
					},
					{
						Name: "test",
						Type: "atest3-t",
					},
				},
			},
		},
		"mixed type len struct should be applied to sorted struct desc": {
			tlen: tlendesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test-1",
					},
					{
						Name: "test",
						Type: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name: "test",
						Type: "test1000000",
					},
					{
						Name:  "test",
						Type:  "zzz",
						Size:  10,
						Align: 10,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test1000000",
					},
					{
						Name: "test",
						Type: "1-test-1",
						Doc:  []string{"test"},
					},
					{
						Name: "test",
						Type: "test-1",
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name:  "test",
						Type:  "zzz",
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
			r, err := tcase.tlen.Apply(tcase.ctx, tcase.o)
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
