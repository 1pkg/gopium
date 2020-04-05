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
	List []StructMock
	Err  error
}

// VisitTop mock implementation
func (w WalkerMock) VisitTop(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	if w.Err != nil {
		return w.Err
	}
	for i := range w.List {
		item := &w.List[i]
		if item.Deep {
			continue
		}
		if r, err := stg.Apply(ctx, item.Struct); err == nil {
			item.Struct = r
		} else {
			return err
		}
	}
	return nil
}

// VisitDeep mock implementation
func (w WalkerMock) VisitDeep(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	if w.Err != nil {
		return w.Err
	}
	for i := range w.List {
		item := &w.List[i]
		if r, err := stg.Apply(ctx, item.Struct); err == nil {
			item.Struct = r
		} else {
			return err
		}
	}
	return nil
}
