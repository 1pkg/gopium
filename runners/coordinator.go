package runners

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"1pkg/gopium"
	"1pkg/gopium/strategies"
)

// coordinator defines helper
// that coordinates runner stages
// - strategy building
// - walker building
// - visiting
type coordinator struct {
	wregex    *regexp.Regexp
	wdeep     bool
	wbackref  bool
	tgroup    string
	tenable   bool
	tforce    bool
	tdiscrete bool
	gtimeout  time.Duration
}

// strategy builds strategy instance
// by using builder and strategies names
func (coord coordinator) strategy(b gopium.StrategyBuilder, snames []gopium.StrategyName) (gopium.Strategy, error) {
	// build strategies
	// and pipe them
	stgs := make([]gopium.Strategy, 0, len(snames))
	for _, sname := range snames {
		stg, err := b.Build(sname)
		if err != nil {
			return nil, fmt.Errorf("can't build such strategy %q %v", sname, err)
		}
		stgs = append(stgs, stg)
	}
	// append tag strategy if enabled
	if coord.tenable {
		tag := strategies.Tag(coord.tgroup, coord.tforce, coord.tdiscrete, snames...)
		stgs = append(stgs, tag)
	}
	return strategies.Pipe(stgs...), nil
}

// walker builds walker instance
// by using builder and walker name
func (coord coordinator) walker(b gopium.WalkerBuilder, wname gopium.WalkerName) (gopium.Walker, error) {
	// build walker
	walker, err := b.Build(wname)
	if err != nil {
		return nil, fmt.Errorf("can't build such walker %q %v", wname, err)
	}
	return walker, nil
}

// visit coordinates walker visiting
func (coord coordinator) visit(ctx context.Context, w gopium.Walker, stg gopium.Strategy) error {
	// set up timeout context
	if coord.gtimeout > 0 {
		nctx, cancel := context.WithTimeout(ctx, coord.gtimeout)
		defer cancel()
		ctx = nctx
	}
	// exec visit on walker with strategy
	if err := w.Visit(
		ctx,
		coord.wregex,
		stg,
		coord.wdeep,
		coord.wbackref,
	); err != nil {
		return fmt.Errorf("strategy error happened %v", err)
	}
	return nil
}
