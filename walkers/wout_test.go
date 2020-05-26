package walkers

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/strategies"
	"1pkg/gopium/tests/data"
	"1pkg/gopium/tests/mocks"
	"1pkg/gopium/typepkg"
)

func TestWout(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	b := strategies.Builder{}
	np, err := b.Build(strategies.Nope)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pck, err := b.Build(strategies.Pack)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 64, 64, 64)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		ctx  context.Context
		r    *regexp.Regexp
		p    gopium.TypeParser
		fmt  gopium.Xbytes
		w    gopium.Writer
		stg  gopium.Strategy
		deep bool
		bref bool
		sts  map[string][]byte
		err  error
	}{
		"empty pkg should visit nothing": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("empty"),
			fmt: mocks.Xbytes{}.Bytes,
			w:   &mocks.Writer{},
			stg: np,
			sts: map[string][]byte{},
		},
		"single struct pkg should visit the struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{}.Bytes,
			w:   &mocks.Writer{},
			stg: np,
			sts: map[string][]byte{
				"tests_data_single_gopium": []byte(`
[
	{
		"Name": "Single",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "A",
				"Type": "string",
				"Size": 16,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "B",
				"Type": "string",
				"Size": 16,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "C",
				"Type": "string",
				"Size": 16,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	}
]
`),
			},
		},
		"single struct pkg should visit nothing on canceled context": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{}.Bytes,
			w:   &mocks.Writer{},
			stg: np,
			sts: map[string][]byte{},
			err: context.Canceled,
		},
		"single struct pkg should visit nothing on parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   mocks.Parser{Typeserr: errors.New("test-1")},
			fmt: mocks.Xbytes{}.Bytes,
			w:   &mocks.Writer{},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-1"),
		},
		"single struct pkg should visit nothing on strategy error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{}.Bytes,
			w:   &mocks.Writer{},
			stg: &mocks.Strategy{Err: errors.New("test-2")},
			sts: map[string][]byte{},
			err: errors.New("test-2"),
		},
		"single struct pkg should visit nothing on writer error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{}.Bytes,
			w:   (&mocks.Writer{Gerr: errors.New("test-3")}),
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-3"),
		},
		"single struct pkg should visit nothing on fmt error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{Err: errors.New("test-4")}.Bytes,
			w:   &mocks.Writer{},
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-4"),
		},
		"single struct pkg should visit nothing on writer persist error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{}.Bytes,
			w: (&mocks.Writer{RWCs: map[string]*mocks.RWC{
				"tests_data_single_gopium": {Werr: errors.New("test-5")},
			}}),
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-5"),
		},
		"single struct pkg should visit nothing on writer close error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Xbytes{}.Bytes,
			w: (&mocks.Writer{RWCs: map[string]*mocks.RWC{
				"tests_data_single_gopium": {Cerr: errors.New("test-6")},
			}}),
			stg: np,
			sts: map[string][]byte{},
			err: errors.New("test-6"),
		},
		"multi structs pkg should visit all expected levels structs with deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`([AZ])`),
			p:    data.NewParser("multi"),
			fmt:  mocks.Xbytes{}.Bytes,
			w:    &mocks.Writer{},
			stg:  pck,
			deep: true,
			sts: map[string][]byte{
				"tests_data_multi_gopium": []byte(`
[
	{
		"Name": "A",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "a",
				"Type": "int64",
				"Size": 8,
				"Align": 8,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	},
	{
		"Name": "AZ",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "D",
				"Type": "1pkg/gopium/tests/data/multi.D",
				"Size": 24,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "a",
				"Type": "bool",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "z",
				"Type": "bool",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	},
	{
		"Name": "Zeze",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "AZ",
				"Type": "1pkg/gopium/tests/data/multi.AZ",
				"Size": 33,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": true,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "D",
				"Type": "1pkg/gopium/tests/data/multi.D",
				"Size": 24,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": true,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "AWA",
				"Type": "1pkg/gopium/tests/data/multi.D",
				"Size": 24,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "ze",
				"Type": "1pkg/gopium/tests/data/multi.ze",
				"Size": 16,
				"Align": 8,
				"Tag": "",
				"Exported": false,
				"Embedded": true,
				"Doc": null,
				"Comment": null
			}
		]
	},
	{
		"Name": "TestAZ",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "D",
				"Type": "1pkg/gopium/tests/data/multi.A",
				"Size": 8,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "a",
				"Type": "bool",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "z",
				"Type": "bool",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	}
]
`),
			},
		},
		"multi structs pkg should visit all expected levels structs without deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`([AZ])`),
			p:    data.NewParser("multi"),
			fmt:  mocks.Xbytes{}.Bytes,
			w:    &mocks.Writer{},
			stg:  pck,
			bref: true,
			sts: map[string][]byte{
				"tests_data_multi_gopium": []byte(`
[
	{
		"Name": "A",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "a",
				"Type": "int64",
				"Size": 8,
				"Align": 8,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	},
	{
		"Name": "AZ",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "D",
				"Type": "1pkg/gopium/tests/data/multi.D",
				"Size": 24,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "a",
				"Type": "bool",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "z",
				"Type": "bool",
				"Size": 1,
				"Align": 1,
				"Tag": "",
				"Exported": false,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			}
		]
	},
	{
		"Name": "Zeze",
		"Doc": null,
		"Comment": null,
		"Fields": [
			{
				"Name": "AZ",
				"Type": "1pkg/gopium/tests/data/multi.AZ",
				"Size": 32,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": true,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "D",
				"Type": "1pkg/gopium/tests/data/multi.D",
				"Size": 24,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": true,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "AWA",
				"Type": "1pkg/gopium/tests/data/multi.D",
				"Size": 24,
				"Align": 8,
				"Tag": "",
				"Exported": true,
				"Embedded": false,
				"Doc": null,
				"Comment": null
			},
			{
				"Name": "ze",
				"Type": "1pkg/gopium/tests/data/multi.ze",
				"Size": 16,
				"Align": 8,
				"Tag": "",
				"Exported": false,
				"Embedded": true,
				"Doc": null,
				"Comment": null
			}
		]
	}
]
`),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			wout := wout{
				fmt:    tcase.fmt,
				writer: tcase.w,
			}.With(tcase.p, m, tcase.deep, tcase.bref)
			// exec
			err := wout.Visit(tcase.ctx, tcase.r, tcase.stg)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			// process checks only on success
			if tcase.err == nil {
				w := tcase.w.(*mocks.Writer)
				for id, rwc := range w.RWCs {
					// check all struct
					// against bytes map
					if st, ok := tcase.sts[id]; ok {
						// read rwc to buffer
						var buf bytes.Buffer
						_, err := buf.ReadFrom(rwc)
						if !reflect.DeepEqual(err, nil) {
							t.Errorf("actual %v doesn't equal to expected %v", err, nil)
						}
						// format actual and expected identically
						actual := strings.Trim(string(buf.Bytes()), "\n")
						expected := strings.Trim(string(st), "\n")
						if !reflect.DeepEqual(actual, expected) {
							t.Errorf("id %v actual %v doesn't equal to expected %v", id, actual, expected)
						}
						delete(tcase.sts, id)
					} else {
						t.Errorf("actual %v doesn't equal to expected %v", id, "")
					}
				}
				// check that map has been drained
				if !reflect.DeepEqual(tcase.sts, map[string][]byte{}) {
					t.Errorf("actual %v doesn't equal to expected %v", tcase.sts, map[string][]byte{})
				}
			}
		})
	}
}
