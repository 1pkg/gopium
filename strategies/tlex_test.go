package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
)

func TestTLexAsc(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself": {
			ctx: context.Background(),
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
			ctx: cctx,
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
		"asc type lex struct should be applied to itself": {
			ctx: context.Background(),
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
			ctx: context.Background(),
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
		"mixed type lex struct should be applied to sorted struct": {
			ctx: context.Background(),
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
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			r, err := tlexasc.Apply(tcase.ctx, tcase.o)
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

func TestTLexDesc(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			ctx: context.Background(),
		},
		"non empty struct should be applied to itself": {
			ctx: context.Background(),
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
			ctx: cctx,
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
		"asc type lex struct should be applied to sorted struct": {
			ctx: context.Background(),
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
			ctx: context.Background(),
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
		"mixed type lex struct should be applied to sorted struct": {
			ctx: context.Background(),
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
			r, err := tlexdesc.Apply(tcase.ctx, tcase.o)
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
