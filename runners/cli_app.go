package runners

import (
	"context"
	"fmt"
	"go/parser"
	"path/filepath"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/pkgs_types"
	"1pkg/gopium/strategy"
	"1pkg/gopium/walker"

	"golang.org/x/tools/go/packages"
)

// CliApp defines cli runner implementation
// that is able to run full gopium cli application
type CliApp struct {
	compiler, arch string
	cpucaches      []int64
	pkg, path      string
	benvs, bflags  []string
	walker         gopium.WalkerName
	regex          *regexp.Regexp
	deep, backref  bool
	strategies     []gopium.StrategyName
}

// NewCliApp helps to spawn new cli application runner
// from list of received parameters
func NewCliApp(
	compiler, arch string,
	cpucaches []int,
	pkg, path string,
	benvs, bflags []string,
	walker, regex string,
	deep, backref bool,
	strategies []string,
) CliApp {
	// cast caches to int64
	caches := make([]int64, 0, len(cpucaches))
	for _, cache := range cpucaches {
		caches = append(caches, int64(cache))
	}
	// cast strategies strings to strategy names
	stgs := make([]gopium.StrategyName, 0, len(strategies))
	for _, strategy := range strategies {
		stgs = append(stgs, gopium.StrategyName(strategy))
	}
	// cast walker string to walker name
	wname := gopium.WalkerName(walker)
	// compile regex
	cregex := regexp.MustCompile(regex)
	// combine cli runner
	return CliApp{
		compiler:   compiler,
		arch:       arch,
		cpucaches:  caches,
		pkg:        pkg,
		path:       path,
		benvs:      benvs,
		bflags:     bflags,
		walker:     wname,
		regex:      cregex,
		deep:       deep,
		backref:    backref,
		strategies: stgs,
	}
}

// Run CliApp implementation
func (cli CliApp) Run(ctx context.Context) error {
	// set up maven
	m := pkgs_types.NewMavenGoTypes(cli.compiler, cli.arch, cli.cpucaches...)
	// set up parser
	absp, err := filepath.Abs(cli.path)
	if err != nil {
		return fmt.Errorf("can't find such path %q %v", cli.path, err)
	}
	p := pkgs_types.ParserXToolPackagesAst{
		Pattern: cli.pkg,
		AbsDir:  absp,
		//nolint
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments,
		BuildEnv:   cli.benvs,
		BuildFlags: cli.bflags,
	}
	// set walker and strategy builders
	wb := walker.NewBuilder(p, m, cli.backref)
	sb := strategy.NewBuilder(m)
	// build strategy
	stgs := make([]gopium.Strategy, 0, len(cli.strategies))
	for _, strategy := range cli.strategies {
		stg, err := sb.Build(gopium.StrategyName(strategy))
		if err != nil {
			return fmt.Errorf("can't build such strategy %q %v", strategy, err)
		}
		stgs = append(stgs, stg)
	}
	stg := strategy.Pipe(stgs...)
	// build walker
	w, err := wb.Build(gopium.WalkerName(cli.walker))
	if err != nil {
		return fmt.Errorf("can't build such walker %q %v", cli.walker, err)
	}
	// run visit function for strategy with regex
	if cli.deep {
		if err = w.VisitDeep(ctx, cli.regex, stg); err != nil {
			return fmt.Errorf("strategy error happened %v", err)
		}
	} else {
		if err = w.VisitTop(ctx, cli.regex, stg); err != nil {
			return fmt.Errorf("strategy error happened %v", err)
		}
	}
	return nil
}
