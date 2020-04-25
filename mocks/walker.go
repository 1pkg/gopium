package mocks

import (
	"context"
	"regexp"

	"1pkg/gopium"
)

// Walker defines mock walker implementation
type Walker struct {
	Structs []gopium.Struct
	Err     error
}

// Visit mock implementation
func (w Walker) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	// check error at start
	if w.Err != nil {
		return w.Err
	}
	// go through structs slice
	for i, st := range w.Structs {
		// apply visiting
		// in case of any error
		// return it back
		// otherwise update struct
		// in the slice
		if r, err := stg.Apply(ctx, st); err == nil {
			w.Structs[i] = r
		} else {
			return err
		}
	}
	return nil
}
