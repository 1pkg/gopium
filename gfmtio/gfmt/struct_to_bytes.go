package gfmt

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"

	"1pkg/gopium"
)

// StructToBytes defines abstraction for
// formatting gopium.Struct to byte slice
type StructToBytes func(gopium.Struct) ([]byte, error)

// PrettyJson defines StructToBytes implementation
// which uses json.Marshal with json.Indent to serialize struct
func PrettyJson(st gopium.Struct) ([]byte, error) {
	// just use json marshal with indent
	return json.MarshalIndent(st, "", "\t")
}

// PrettyXml defines StructToBytes implementation
// which uses xml.MarshalIndent to serialize struct
func PrettyXml(st gopium.Struct) ([]byte, error) {
	// just use xml marshal with indent
	return xml.MarshalIndent(st, "", "\t")
}

// PrettyCsv defines StructToBytes implementation
// that serializes struct to csv format
func PrettyCsv(st gopium.Struct) ([]byte, error) {
	// prepare buf and csv writer
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	// write header
	if err := w.Write([]string{
		"Struct Name",
		"Struct Doc",
		"Struct Comment",
		"Field Name",
		"Field Type",
		"Field Size",
		"Field Align",
		"Field Tag",
		"Field Exported",
		"Field Embedded",
		"Field Doc",
		"Field Comment",
	}); err != nil {
		// this should never happen/covered
		return nil, err
	}
	// go through all fields
	// and write then one by one
	for _, f := range st.Fields {
		if err := w.Write([]string{
			st.Name,
			strings.Join(st.Doc, " "),
			strings.Join(st.Comment, " "),
			f.Name,
			f.Type,
			strconv.Itoa(int(f.Size)),
			strconv.Itoa(int(f.Align)),
			f.Tag,
			strconv.FormatBool(f.Exported),
			strconv.FormatBool(f.Embedded),
			strings.Join(f.Doc, " "),
			strings.Join(f.Comment, " "),
		}); err != nil {
			// this should never happen/covered
			return nil, err
		}
	}
	// flush to buf
	w.Flush()
	// check flush error
	if err := w.Error(); err != nil {
		// this should never happen/covered
		return nil, err
	}
	// and return buf result
	return buf.Bytes(), nil
}
