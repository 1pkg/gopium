package mocks

import (
	"encoding/json"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// Diff defines mock fmtio diff implementation
type Diff struct {
	Err error
}

// Diff mock implementation
func (fmt Diff) Diff(ho collections.Hierarchic, hr collections.Hierarchic) ([]byte, error) {
	// in case we have error
	// return it back
	if fmt.Err != nil {
		return nil, fmt.Err
	}
	// otherwise use json bytes impl
	data := [][]gopium.Struct{ho.Flat().Sorted(), hr.Flat().Sorted()}
	return json.MarshalIndent(data, "", "\t")
}
