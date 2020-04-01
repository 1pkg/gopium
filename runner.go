package gopium

import "context"

// Runner defines abstraction for
// simple root gopium runner
type Runner interface {
	Run(context.Context) error
}
