package runners

import (
	"context"
	"fmt"
	"go/build"
	"go/parser"
	"regexp"
	"strings"
	"time"

	"1pkg/gopium"
	"1pkg/gopium/fmtio"
	"1pkg/gopium/strategies"
	"1pkg/gopium/typepkg"
	"1pkg/gopium/walkers"

	"golang.org/x/tools/go/packages"
)

// CliApp defines cli runner implementation
// that is able to run full gopium cli application
type CliApp struct {
	coord    coordinator
	wbuilder gopium.WalkerBuilder
	sbuilder gopium.StrategyBuilder
	wname    gopium.WalkerName
	snames   []gopium.StrategyName
}

// NewCliApp helps to spawn new cli application runner
// from list of received parameters or returns error
func NewCliApp(
	// target platform vars
	compiler,
	arch string,
	cpucaches []int,
	// package parser vars
	pkg,
	path string,
	benvs,
	bflags []string,
	// walker vars
	walker,
	regex string,
	deep,
	backref bool,
	stgs []string,
	// printer vars
	indent,
	tabwidth int,
	usespace bool,
	// global vars
	timeout int,
) (*CliApp, error) {
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
	// set up parser
	p := typepkg.ParserXToolPackagesAst{
		Pattern:    pkg,
		Root:       build.Default.GOPATH,
		Path:       path,
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments | parser.AllErrors,
		BuildEnv:   benvs,
		BuildFlags: bflags,
	}
	// set up printer
	pr := fmtio.Goprint(indent, tabwidth, usespace)
	// compile regexp
	wregex, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("can't compile such regexp %q %v", wregex, err)
	}
	// cast timeout to second duration
	gtimeout := time.Duration(timeout) * time.Second
	// set up coordinator
	coord := coordinator{
		wregex:   wregex,
		gtimeout: gtimeout,
	}
	// set walker and strategy builders
	wbuilder := walkers.Builder{
		Parser:  p,
		Exposer: m,
		Printer: pr,
		Deep:    deep,
		Bref:    backref,
	}
	sbuilder := strategies.Builder{Curator: m}
	// cast strategies strings to strategy names
	snames := make([]gopium.StrategyName, 0, len(stgs))
	for _, strategy := range stgs {
		snames = append(snames, gopium.StrategyName(strategy))
	}
	// cast walker string to walker name
	wname := gopium.WalkerName(walker)
	// combine cli runner
	return &CliApp{
		coord:    coord,
		wbuilder: wbuilder,
		sbuilder: sbuilder,
		wname:    wname,
		snames:   snames,
	}, nil
}

// Run CliApp implementation
func (cli *CliApp) Run(ctx context.Context) error {
	// build strategy
	stg, err := cli.coord.strategy(cli.sbuilder, cli.snames)
	if err != nil {
		return err
	}
	// build walker
	w, err := cli.coord.walker(cli.wbuilder, cli.wname)
	if err != nil {
		return err
	}
	// run visit
	return cli.coord.visit(ctx, w, stg)
}
