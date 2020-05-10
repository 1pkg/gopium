package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestTLex(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		tlex tlex
		ctx  context.Context
		o    gopium.Struct
		r    gopium.Struct
		err  error
	}{
		"empty struct should be applied to empty struct": {
			tlex: tlexasc,
			ctx:  context.Background(),
		},
		"non empty struct should be applied to itself": {
			tlex: tlexasc,
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
			tlex: tlexdesc,
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
			err: context.Canceled,
		},
		"asc type lex struct should be applied to itself": {
			tlex: tlexasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test1",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test3",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test1",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test3",
					},
				},
			},
		},
		"desc type lex struct should be applied to sorted struct": {
			tlex: tlexasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test3",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test1",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test1",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test3",
					},
				},
			},
		},
		"mixed type lex struct should be applied to sorted struct asc": {
			tlex: tlexasc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "ztest1",
					},
					{
						Name: "test",
						Type: "test2",
						Doc:  []string{"test"},
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name: "test",
						Type: "test4",
					},
					{
						Name:  "test",
						Type:  "ttest5",
						Size:  10,
						Align: 10,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name: "test",
						Type: "test2",
						Doc:  []string{"test"},
					},
					{
						Name: "test",
						Type: "test4",
					},
					{
						Name:  "test",
						Type:  "ttest5",
						Size:  10,
						Align: 10,
					},
					{
						Name: "test",
						Type: "ztest1",
					},
				},
			},
		},
		"asc type lex struct should be applied to sorted struct": {
			tlex: tlexdesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test1",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test3",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test3",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test1",
					},
				},
			},
		},
		"desc type lex struct should be applied to itself": {
			tlex: tlexdesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test3",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test1",
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "test3",
					},
					{
						Name: "test",
						Type: "test2",
					},
					{
						Name: "test",
						Type: "test1",
					},
				},
			},
		},
		"mixed type lex struct should be applied to sorted struct desc": {
			tlex: tlexdesc,
			ctx:  context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Type: "ztest1",
					},
					{
						Name: "test",
						Type: "test2",
						Doc:  []string{"test"},
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
					{
						Name: "test",
						Type: "test4",
					},
					{
						Name:  "test",
						Type:  "ttest5",
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
						Type: "ztest1",
					},
					{
						Name:  "test",
						Type:  "ttest5",
						Size:  10,
						Align: 10,
					},
					{
						Name: "test",
						Type: "test4",
					},
					{
						Name: "test",
						Type: "test2",
						Doc:  []string{"test"},
					},
					{
						Name:     "test",
						Type:     "atest3",
						Exported: true,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tcase.tlex.Apply(tcase.ctx, tcase.o)
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
