package runners

import (
	"context"
	"fmt"
	"go/parser"
	"path/filepath"
	"regexp"
	"time"

	"1pkg/gopium"
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
	timeout  time.Duration
}

// NewCliApp helps to spawn new cli application runner
// from list of received parameters or returns error
func NewCliApp(
	compiler, arch string,
	cpucaches []int,
	pkg, path string,
	benvs, bflags []string,
	walker, regex string,
	deep, backref bool,
	stgs []string,
	tgroup string,
	tenable, tforce, tdiscrete bool,
	timeout int,
) (*CliApp, error) {
	// cast caches to int64
	caches := make([]int64, 0, len(cpucaches))
	for _, cache := range cpucaches {
		caches = append(caches, int64(cache))
	}
	// set up maven
	m := typepkg.NewMavenGoTypes(compiler, arch, caches...)
	// set up parser
	absp, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("can't find such path %q %v", path, err)
	}
	p := typepkg.ParserXToolPackagesAst{
		Pattern: pkg,
		AbsDir:  absp,
		//nolint
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments,
		BuildEnv:   benvs,
		BuildFlags: bflags,
	}
	// compile regexp
	wregex, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("can't compile such regexp %q %v", wregex, err)
	}
	// set up coordinator
	coord := coordinator{
		wregex:    wregex,
		wdeep:     deep,
		tgroup:    tgroup,
		tenable:   tenable,
		tforce:    tforce,
		tdiscrete: tdiscrete,
	}
	// set walker and strategy builders
	wbuilder := walkers.NewBuilder(p, m, backref)
	sbuilder := strategies.NewBuilder(m)
	// cast strategies strings to strategy names
	snames := make([]gopium.StrategyName, 0, len(stgs))
	for _, strategy := range stgs {
		snames = append(snames, gopium.StrategyName(strategy))
	}
	// cast walker string to walker name
	wname := gopium.WalkerName(walker)
	// cast timeout to second duration
	tm := time.Duration(timeout) * time.Second
	// combine cli runner
	return &CliApp{
		coord:    coord,
		wbuilder: wbuilder,
		sbuilder: sbuilder,
		wname:    wname,
		snames:   snames,
		timeout:  tm,
	}, nil
}

// Run CliApp implementation
func (cli CliApp) Run(ctx context.Context) error {
	// set up timeout context
	if cli.timeout > 0 {
		nctx, cancel := context.WithTimeout(ctx, cli.timeout)
		defer cancel()
		ctx = nctx
	}
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
