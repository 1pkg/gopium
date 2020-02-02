package fmts

import (
	"bytes"
	"encoding/json"
)

// TypeFormat defines abstraction for
// formatting generic type to byte slice
type TypeFormat func(interface{}) ([]byte, error)

// PrettyJson defines json.Marshal with json.Indent TypeFormat implementation
func PrettyJson(i interface{}) ([]byte, error) {
	r, err := json.Marshal(i)
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
