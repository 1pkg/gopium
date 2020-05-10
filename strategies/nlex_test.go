package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestNLex(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		nlex nlex
		ctx  context.Context
		o    gopium.Struct
		r    gopium.Struct
		err  error
	}{
		"empty struct should be applied to empty struct": {
			nlex: nlexasc,
			ctx:  context.Background(),
		},
		"non empty struct should be applied to itself": {
			nlex: nlexasc,
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
			nlex: nlexdesc,
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
			err: context.Canceled,
		},
		"asc name lex struct should be applied to itself": {
			nlex: nlexasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
					{
						Name: "test3",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
					{
						Name: "test3",
					},
				},
			},
		},
		"desc name lex struct should be applied to sorted struct": {
			nlex: nlexasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test3",
					},
					{
						Name: "test2",
					},
					{
						Name: "test1",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
					{
						Name: "test3",
					},
				},
			},
		},
		"mixed name lex struct should be applied to sorted struct asc": {
			nlex: nlexasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "ztest1",
					},
					{
						Name: "test2",
						Doc:  []string{"test"},
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name: "test4",
					},
					{
						Name:  "ttest5",
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
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name: "test2",
						Doc:  []string{"test"},
					},
					{
						Name: "test4",
					},
					{
						Name:  "ttest5",
						Type:  "int33",
						Size:  10,
						Align: 10,
					},
					{
						Name: "ztest1",
					},
				},
			},
		},
		"asc name lex struct should be applied to sorted struct": {
			nlex: nlexdesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
					},
					{
						Name: "test2",
					},
					{
						Name: "test3",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test3",
					},
					{
						Name: "test2",
					},
					{
						Name: "test1",
					},
				},
			},
		},
		"desc name lex struct should be applied to itself": {
			nlex: nlexdesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test3",
					},
					{
						Name: "test2",
					},
					{
						Name: "test1",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test3",
					},
					{
						Name: "test2",
					},
					{
						Name: "test1",
					},
				},
			},
		},
		"mixed name lex struct should be applied to sorted struct desc": {
			nlex: nlexdesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "ztest1",
					},
					{
						Name: "test2",
						Doc:  []string{"test"},
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
					{
						Name: "test4",
					},
					{
						Name:  "ttest5",
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
						Name: "ztest1",
					},
					{
						Name:  "ttest5",
						Type:  "int33",
						Size:  10,
						Align: 10,
					},
					{
						Name: "test4",
					},
					{
						Name: "test2",
						Doc:  []string{"test"},
					},
					{
						Name:     "atest3",
						Type:     "int64",
						Exported: true,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.nlex.Apply(tcase.ctx, tcase.o)
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
