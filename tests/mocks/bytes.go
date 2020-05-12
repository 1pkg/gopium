package mocks

import (
	"encoding/json"

	"1pkg/gopium/collections"
)

// Bytes defines mock fmtio bytes implementation
type Bytes struct {
	Err error
}

// Bytes mock implementation
func (fmt Bytes) Bytes(f collections.Flat) ([]byte, error) {
	// in case we have error
	// return it back
	if fmt.Err != nil {
		return nil, fmt.Err
	}
	// otherwise use json bytes impl
	return json.MarshalIndent(f.Sorted(), "", "\t")
}
