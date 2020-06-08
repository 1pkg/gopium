package gopium

import "context"

// Runner defines abstraction for gopium runner
type Runner interface {
	Run(context.Context) error
}
