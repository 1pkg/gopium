package mocks

import (
	"context"
	"regexp"

	"1pkg/gopium"
)

// StructMock defines mock gopium struct
// data transfer object with deep flag
type StructMock struct {
	gopium.Struct
	Deep bool
}

// WalkerMock defines mock walker implementation
type WalkerMock struct {
	Structs []StructMock
	Err     error
}

// VisitTop mock implementation
func (w WalkerMock) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy, deep bool) error {
	// check error at start
	if w.Err != nil {
		return w.Err
	}
	// go through structs slice
	for i := range w.Structs {
		st := &w.Structs[i]
		// in case type of visiting
		// doesn't match skip it
		if st.Deep != deep {
			continue
		}
		// otherwise apply visiting
		// in case of any error
		// return it back
		// otherwise update struct in the slice
		if r, err := stg.Apply(ctx, st.Struct); err == nil {
			st.Struct = r
		} else {
			return err
		}
	}
	return nil
}
