package mocks

import (
	"1pkg/gopium"
	"1pkg/gopium/fmtio"
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
	return fmtio.Json(st)
}
