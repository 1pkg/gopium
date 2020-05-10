package ctx

import (
	"context"
	"os"
	"os/signal"
)

// WithSignals creates context with cancelation
// which listens to provided list of signals
func WithSignals(ctx context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	// prepare global context
	// with cancelation
	// on system signals
	ctx, cancel := context.WithCancel(ctx)
	// run separate listener goroutine
	go func() {
		defer cancel()
		// prepare signal chan for
		// global context cancelation
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, sigs...)
		// on signal or cancelation
		// stop the goroutine
		select {
		case <-ctx.Done():
		case <-sig:
		}
	}()
	return ctx, cancel
}