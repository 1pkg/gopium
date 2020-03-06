package fmts

import (
	"bytes"
	"encoding/json"

	"1pkg/gopium"
)

// StructToBytes defines abstraction for
// formatting gopium.Struct to byte slice
type StructToBytes func(gopium.Struct) ([]byte, error)

// PrettyJson defines json.Marshal
// with json.Indent StructToBytes implementation
func PrettyJson(st gopium.Struct) ([]byte, error) {
	r, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, r, "", "\t")
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
