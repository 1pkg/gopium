package mocks

import "context"

// Runner defines mock runner implementation
type Runner struct {
	Err error
}

// Run mock implementation
func (r Runner) Run(context.Context) error {
	return r.Err
}
