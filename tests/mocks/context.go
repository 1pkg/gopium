package mocks

import (
	"context"
	"time"
)

// Context defines mock context implementation
type Context struct {
	After   int
	current int
}

// Deadline mock implementation
func (ctx Context) Deadline() (deadline time.Time, ok bool) {
	return context.Background().Deadline()
}

// Done mock implementation
func (ctx *Context) Done() <-chan struct{} {
	ctx.current++
	ch := make(chan struct{})
	// after n-th call to done
	// close the chan
	if ctx.current >= ctx.After {
		close(ch)
	}
	return ch
}

// Err mock implementation
func (ctx Context) Err() error {
	// after n-th call to done
	// return ctx error
	if ctx.current >= ctx.After {
		return context.Canceled
	}
	return nil
}

// Value mock implementation
func (ctx Context) Value(key interface{}) interface{} {
	return context.Background().Value(key)
}
