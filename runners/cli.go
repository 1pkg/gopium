package runners

import (
	"context"
	"fmt"
	"go/build"
	"go/parser"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/1pkg/gopium/fmtio"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/strategies"
	"github.com/1pkg/gopium/typepkg"
	"github.com/1pkg/gopium/walkers"

	"golang.org/x/tools/go/packages"
)

// Cli defines cli runner implementation
// that is able to run full gopium cli application
type Cli struct {
	snames []gopium.StrategyName  `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	v      visitor                `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	wb     gopium.WalkerBuilder   `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	sb     gopium.StrategyBuilder `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	wname  gopium.WalkerName      `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_      [40]byte               `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 128 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; - 🌺 gopium @1pkg

// NewCli helps to spawn new cli application runner
// from list of received parameters or returns error
func NewCli(
	// target platform vars
	compiler,
	arch string,
	cpucaches []int,
	// package parser vars
	pkg,
	path string,
	benvs,
	bflags []string,
	// gopium walker vars
	walker,
	regex string,
	deep,
	backref bool,
	stgs []string,
	// gopium printer vars
	indent,
	tabwidth int,
	usespace,
	usegofmt bool,
	// gopium global vars
	timeout int,
) (*Cli, error) {
	// cast caches to int64
	caches := make([]int64, 0, len(cpucaches))
	for _, cache := range cpucaches {
		caches = append(caches, int64(cache))
	}
	// set up maven
	m, err := typepkg.NewMavenGoTypes(compiler, arch, caches...)
	if err != nil {
		return nil, fmt.Errorf("can't set up maven %v", err)
	}
	// replace package template
	path = strings.Replace(path, "{{package}}", pkg, 1)
	// set root to gopath only if
	// not absolute path has been provided
	var root string
	if !filepath.IsAbs(path) {
		root = build.Default.GOPATH
	}
	// https://github.com/1pkg/gopium/issues/18
	if len(benvs) == 0 {
		benvs = []string{
			fmt.Sprintf("GOPATH=%s", build.Default.GOPATH),
			fmt.Sprintf("GOCACHE=%s", filepath.Join(build.Default.GOPATH, ".cache")),
		}
	}
	// set up parser
	xp := &typepkg.ParserXToolPackagesAst{
		Pattern: pkg,
		Root:    root,
		Path:    path,
		//nolint
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments | parser.AllErrors,
		BuildEnv:   benvs,
		BuildFlags: bflags,
	}
	// set up printer
	var p gopium.Printer
	if usegofmt {
		p = fmtio.Gofmt{}
	} else {
		p = fmtio.NewGoprinter(indent, tabwidth, usespace)
	}
	// compile regexp
	cregex, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("can't compile such regexp %v", err)
	}
	// cast timeout to second duration
	stimeout := time.Duration(timeout) * time.Second
	// set up visitor
	v := visitor{
		regex:   cregex,
		timeout: stimeout,
	}
	// set walker and strategy builders
	wb := walkers.Builder{
		Parser:  xp,
		Exposer: m,
		Printer: p,
		Deep:    deep,
		Bref:    backref,
	}
	sb := strategies.Builder{Curator: m}
	// cast strategies strings to strategy names
	snames := make([]gopium.StrategyName, 0, len(stgs))
	for _, strategy := range stgs {
		snames = append(snames, gopium.StrategyName(strategy))
	}
	// cast walker string to walker name
	wname := gopium.WalkerName(walker)
	// combine cli runner
	return &Cli{
		v:      v,
		wb:     wb,
		sb:     sb,
		wname:  wname,
		snames: snames,
	}, nil
}

// Run cli implementation
func (cli *Cli) Run(ctx context.Context) error {
	// build strategy
	stg, err := cli.v.strategy(cli.sb, cli.snames)
	if err != nil {
		return err
	}
	// build walker
	w, err := cli.v.walker(cli.wb, cli.wname)
	if err != nil {
		return err
	}
	// run visitor visiting
	return cli.v.visit(ctx, w, stg)
}
