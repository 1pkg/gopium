package runners

import (
	"context"
	"errors"
	"go/build"
	"go/parser"
	"reflect"
	"regexp"
	"testing"
	"time"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
	"1pkg/gopium/strategies"
	"1pkg/gopium/tests/mocks"
	"1pkg/gopium/typepkg"
	"1pkg/gopium/walkers"

	"golang.org/x/tools/go/packages"
)

func TestNewCli(t *testing.T) {
	// prepare
	m, err := typepkg.NewMavenGoTypes("gc", "amd64", 2, 4, 8)
	if !reflect.DeepEqual(err, nil) {
		t.Fatalf("actual %v doesn't equal to expected %v", err, nil)
	}
	table := map[string]struct {
		// target platform vars
		compiler  string
		arch      string
		cpucaches []int
		// package parser vars
		pkg    string
		path   string
		benvs  []string
		bflags []string
		// walker vars
		walker  string
		regex   string
		deep    bool
		backref bool
		stgs    []string
		// printer vars
		indent   int
		tabwidth int
		usespace bool
		// global vars
		timeout int
		// test vars
		cli *Cli
		err error
	}{
		"new cli cli should expected cli on valid parameters": {
			// target platform vars
			compiler:  "gc",
			arch:      "amd64",
			cpucaches: []int{2, 4, 8},
			// package parser vars
			pkg:    "test-pkg",
			path:   "test-path",
			benvs:  []string{},
			bflags: []string{},
			// walker vars
			walker:  "test-w",
			regex:   `.*`,
			deep:    true,
			backref: true,
			stgs:    []string{"test-stg"},
			// printer vars
			indent:   4,
			tabwidth: 4,
			usespace: true,
			// global vars
			timeout: 5,
			// test vars
			cli: &Cli{
				v: visitor{
					regex:   regexp.MustCompile(`.*`),
					timeout: 5 * time.Second,
				},
				wb: walkers.Builder{
					Parser: &typepkg.ParserXToolPackagesAst{
						Pattern:    "test-pkg",
						Root:       build.Default.GOPATH,
						Path:       "test-path",
						ModeTypes:  packages.LoadAllSyntax,
						ModeAst:    parser.ParseComments | parser.AllErrors,
						BuildEnv:   []string{},
						BuildFlags: []string{},
					},
					Exposer: m,
					Printer: fmtio.NewGoprinter(4, 4, true),
					Deep:    true,
					Bref:    true,
				},
				sb:     strategies.Builder{Curator: m},
				wname:  "test-w",
				snames: []gopium.StrategyName{"test-stg"},
			},
		},
		"new cli cli should expected cli on valid parameters with abs path": {
			// target platform vars
			compiler:  "gc",
			arch:      "amd64",
			cpucaches: []int{2, 4, 8},
			// package parser vars
			pkg:    "test-pkg",
			path:   "/test-path",
			benvs:  []string{},
			bflags: []string{},
			// walker vars
			walker:  "test-w",
			regex:   `.*`,
			deep:    true,
			backref: true,
			stgs:    []string{"test-stg"},
			// printer vars
			indent:   4,
			tabwidth: 4,
			usespace: true,
			// global vars
			timeout: 5,
			// test vars
			cli: &Cli{
				v: visitor{
					regex:   regexp.MustCompile(`.*`),
					timeout: 5 * time.Second,
				},
				wb: walkers.Builder{
					Parser: &typepkg.ParserXToolPackagesAst{
						Pattern:    "test-pkg",
						Path:       "/test-path",
						ModeTypes:  packages.LoadAllSyntax,
						ModeAst:    parser.ParseComments | parser.AllErrors,
						BuildEnv:   []string{},
						BuildFlags: []string{},
					},
					Exposer: m,
					Printer: fmtio.NewGoprinter(4, 4, true),
					Deep:    true,
					Bref:    true,
				},
				sb:     strategies.Builder{Curator: m},
				wname:  "test-w",
				snames: []gopium.StrategyName{"test-stg"},
			},
		},
		"new cli should return error on invalid compiler arch combination": {
			// target platform vars
			compiler:  "cg",
			arch:      "64amd64",
			cpucaches: []int{2, 4, 8},
			// package parser vars
			pkg:    "test-pkg",
			path:   "test-path",
			benvs:  []string{},
			bflags: []string{},
			// walker vars
			walker:  "test-w",
			regex:   `.*`,
			deep:    true,
			backref: true,
			stgs:    []string{"test-stg"},
			// printer vars
			indent:   4,
			tabwidth: 4,
			usespace: true,
			// global vars
			timeout: 5,
			// test vars
			err: errors.New(`can't set up maven unsuported compiler "cg" arch "64amd64" combination`),
		},
		"new cli should return error on regex compile error": {
			// target platform vars
			compiler:  "gc",
			arch:      "amd64",
			cpucaches: []int{2, 4, 8},
			// package parser vars
			pkg:    "test-pkg",
			path:   "test-path",
			benvs:  []string{},
			bflags: []string{},
			// walker vars
			walker:  "test-w",
			regex:   `[`,
			deep:    true,
			backref: true,
			stgs:    []string{"test-stg"},
			// printer vars
			indent:   4,
			tabwidth: 4,
			usespace: true,
			// global vars
			timeout: 5,
			// test vars
			err: errors.New("can't compile such regexp error parsing regexp: missing closing ]: `[`"),
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			cli, err := NewCli(
				tcase.compiler,
				tcase.arch,
				tcase.cpucaches,
				tcase.pkg,
				tcase.path,
				tcase.benvs,
				tcase.bflags,
				tcase.walker,
				tcase.regex,
				tcase.deep,
				tcase.backref,
				tcase.stgs,
				tcase.indent,
				tcase.tabwidth,
				tcase.usespace,
				tcase.timeout,
			)
			// check
			if !reflect.DeepEqual(cli, tcase.cli) {
				t.Errorf("actual %v doesn't equal to expected %v", cli, tcase.cli)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}

func TestCliRun(t *testing.T) {
	// prepare
	table := map[string]struct {
		cli *Cli
		err error
	}{
		"cli should return error on strategy builder error": {
			cli: &Cli{
				v:  visitor{},
				sb: mocks.StrategyBuilder{Err: errors.New("test-1")},
			},
			err: errors.New("can't build such strategy [] test-1"),
		},
		"cli should return error on walker builder error": {
			cli: &Cli{
				v:  visitor{},
				sb: mocks.StrategyBuilder{Strategy: &mocks.Strategy{}},
				wb: mocks.WalkerBuilder{Err: errors.New("test-2")},
			},
			err: errors.New(`can't build such walker "" test-2`),
		},
		"cli should return error on visiting error": {
			cli: &Cli{
				v:  visitor{},
				sb: mocks.StrategyBuilder{Strategy: &mocks.Strategy{}},
				wb: mocks.WalkerBuilder{Walker: mocks.Walker{Err: errors.New("test-3")}},
			},
			err: errors.New("visiting error happened test-3"),
		},
		"cli should return error on timeout": {
			cli: &Cli{
				v:  visitor{timeout: time.Nanosecond},
				sb: mocks.StrategyBuilder{Strategy: &mocks.Strategy{}},
				wb: mocks.WalkerBuilder{Walker: mocks.Walker{}},
			},
			err: errors.New("visiting error happened context deadline exceeded"),
		},
		"cli should return expected results on visiting": {
			cli: &Cli{
				v:  visitor{},
				sb: mocks.StrategyBuilder{Strategy: &mocks.Strategy{}},
				wb: mocks.WalkerBuilder{Walker: mocks.Walker{}},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			err := tcase.cli.Run(context.Background())
			// check
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
