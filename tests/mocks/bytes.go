package mocks

import (
	"encoding/json"

	"1pkg/gopium"
)

// Bytes defines mock fmtio bytes implementation
type Bytes struct {
	Err error
}

// Bytes mock implementation
func (fmt Bytes) Bytes(st gopium.Struct) ([]byte, error) {
	// in case we have error
	// return it back
	if fmt.Err != nil {
		return nil, fmt.Err
	}
	// otherwise use json bytes impl
	return json.MarshalIndent(st, "", "\t")
}
