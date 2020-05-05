package walkers

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
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
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	pck, err := b.Build(strategies.Pack)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 64, 64, 64)
	if err != nil {
		t.Fatalf("actual %v doesn't equal to %v", err, nil)
	}
	table := map[string]struct {
		ctx  context.Context
		r    *regexp.Regexp
		p    gopium.TypeParser
		fmt  fmtio.Bytes
		w    fmtio.Writer
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
			fmt: mocks.Bytes{}.Bytes,
			stg: np,
			sts: make(map[string][]byte),
		},
		"single struct pkg should visit the single struct": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{}.Bytes,
			stg: np,
			sts: map[string][]byte{
				"5-b0652be9c761c2f34deff8a560333dd372ee062bb1dbcba6a79647fdc3205919": []byte(`
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
`),
			},
		},
		"single struct pkg should visit nothing on context cancelation": {
			ctx: cctx,
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{}.Bytes,
			stg: np,
			sts: make(map[string][]byte),
			err: cctx.Err(),
		},
		"single struct pkg should visit nothing on parser error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   mocks.Parser{Terr: errors.New("test-1")},
			fmt: mocks.Bytes{}.Bytes,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-1"),
		},
		"single struct pkg should visit nothing on strategy error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{}.Bytes,
			stg: mocks.Strategy{Err: errors.New("test-2")},
			sts: make(map[string][]byte),
			err: errors.New("test-2"),
		},
		"single struct pkg should visit nothing on writer gen error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{}.Bytes,
			w:   (&mocks.Writer{Err: errors.New("test-3")}).Writer,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-3"),
		},
		"single struct pkg should visit nothing on fmt error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{Err: errors.New("test-4")}.Bytes,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-4"),
		},
		"single struct pkg should visit nothing on writer write error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{}.Bytes,
			w:   (&mocks.Writer{Werr: errors.New("test-5")}).Writer,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-5"),
		},
		"single struct pkg should visit nothing on writer close error": {
			ctx: context.Background(),
			r:   regexp.MustCompile(`.*`),
			p:   data.NewParser("single"),
			fmt: mocks.Bytes{}.Bytes,
			w:   (&mocks.Writer{Cerr: errors.New("test-6")}).Writer,
			stg: np,
			sts: make(map[string][]byte),
			err: errors.New("test-6"),
		},
		"multi structs pkg should visit all relevant levels structs with deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`(A|Z)`),
			p:    data.NewParser("multi"),
			fmt:  mocks.Bytes{}.Bytes,
			stg:  pck,
			deep: true,
			sts: map[string][]byte{
				"9-7d858286ee3f6bdbb9c740b5333435af40ec918bdeec00ececacf5ab9764f09b": []byte(`
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
}
`),
				"17-342e1133d9f044ad74cd048f681aad0efcca3407b8fe3b972c96eb92d034fd04": []byte(`
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
}
`),
				"27-6a3c1ba2a278b9b24c0d76ad232bba0f0b0abd806f9cbb6e0910966f761e5130": []byte(`
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
}
`),
				"29-6dc854454cff4b7c6b7ba90ba55fa564c21409c5a107cf402dd2e582d44dd32a": []byte(`
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
`),
			},
		},
		"multi structs pkg should visit all relevant levels structs without deep": {
			ctx:  context.Background(),
			r:    regexp.MustCompile(`(A|Z)`),
			p:    data.NewParser("multi"),
			fmt:  mocks.Bytes{}.Bytes,
			stg:  pck,
			bref: true,
			sts: map[string][]byte{
				"9-7d858286ee3f6bdbb9c740b5333435af40ec918bdeec00ececacf5ab9764f09b": []byte(`
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
}
`),
				"17-342e1133d9f044ad74cd048f681aad0efcca3407b8fe3b972c96eb92d034fd04": []byte(`
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
}
`),
				"27-6a3c1ba2a278b9b24c0d76ad232bba0f0b0abd806f9cbb6e0910966f761e5130": []byte(`
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
`),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			w := &mocks.Writer{}
			wout := wout{
				fmt:    tcase.fmt,
				writer: w.Writer,
			}.With(tcase.p, m, tcase.deep, tcase.bref)
			if tcase.w != nil {
				wout.writer = tcase.w
			}
			err := wout.Visit(tcase.ctx, tcase.r, tcase.stg)
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
			for id, buf := range w.Buffers {
				// check all struct
				// against bytes map
				if st, ok := tcase.sts[id]; ok {
					// format actual and expected identically
					stract, strexp := strings.Trim(string(buf.Bytes()), "\n"), strings.Trim(string(st), "\n")
					if !reflect.DeepEqual(stract, strexp) {
						t.Errorf("id %v actual %v doesn't equal to expected %v", id, stract, strexp)
					}
					delete(tcase.sts, id)
				} else {
					t.Errorf("actual %v doesn't equal to expected %v", id, "")
				}
			}
			// check that map has been drained
			if !reflect.DeepEqual(tcase.sts, make(map[string][]byte)) {
				t.Errorf("actual %v doesn't equal to expected %v", tcase.sts, make(map[string][]byte))
			}
		})
	}
}
