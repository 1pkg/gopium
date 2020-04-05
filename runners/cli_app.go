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
	compiler, arch string
	cpucaches      []int64
	pkg, path      string
	benvs, bflags  []string
	walker         gopium.WalkerName
	regex          *regexp.Regexp
	deep, backref  bool
	strategies     []gopium.StrategyName
	tagtype        TagType
	timeout        time.Duration
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
	tagtype string,
	timeout int,
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
	// cast tagtype string to tag type
	tt := TagType(tagtype)
	// cast timeout to second duration
	tm := time.Duration(timeout) * time.Second
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
		tagtype:    tt,
		timeout:    tm,
	}
}

// Run CliApp implementation
func (cli CliApp) Run(ctx context.Context) error {
	// set up timeout context
	if cli.timeout > 0 {
		nctx, cancel := context.WithTimeout(ctx, cli.timeout)
		defer cancel()
		ctx = nctx
	}
	// set up maven
	m := typepkg.NewMavenGoTypes(cli.compiler, cli.arch, cli.cpucaches...)
	// set up parser
	absp, err := filepath.Abs(cli.path)
	if err != nil {
		return fmt.Errorf("can't find such path %q %v", cli.path, err)
	}
	p := typepkg.ParserXToolPackagesAst{
		Pattern: cli.pkg,
		AbsDir:  absp,
		//nolint
		ModeTypes:  packages.LoadAllSyntax,
		ModeAst:    parser.ParseComments,
		BuildEnv:   cli.benvs,
		BuildFlags: cli.bflags,
	}
	// set walker and strategy builders
	wb := walkers.NewBuilder(p, m, cli.backref)
	sb := strategies.NewBuilder(m)
	// build strategy
	stgs := make([]gopium.Strategy, 0, len(cli.strategies))
	for _, strategy := range cli.strategies {
		stg, err := sb.Build(gopium.StrategyName(strategy))
		if err != nil {
			return fmt.Errorf("can't build such strategy %q %v", strategy, err)
		}
		stgs = append(stgs, stg)
	}
	// append tag strategy
	if cli.tagtype != None {
		force := cli.tagtype == Force
		tag := strategies.Tag(force, cli.strategies...)
		stgs = append(stgs, tag)
	}
	stg := strategies.Pipe(stgs...)
	// build walker
	w, err := wb.Build(gopium.WalkerName(cli.walker))
	if err != nil {
		return fmt.Errorf("can't build such walker %q %v", cli.walker, err)
	}
	// run visit function for strategy with regex
	if err = w.Visit(ctx, cli.regex, stg, cli.deep); err != nil {
		return fmt.Errorf("strategy error happened %v", err)
	}
	return nil
}
