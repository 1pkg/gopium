package gopium

import (
	"bytes"
	"encoding/json"
)

// Formatter defines abstraction for
// formatting generic type to byte slice
type Formatter func(interface{}) ([]byte, error)

// PrettyJson defines json marshal with tab indent Formatter implementation
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
