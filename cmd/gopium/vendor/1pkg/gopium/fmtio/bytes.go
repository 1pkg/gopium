package fmtio

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"1pkg/gopium"
)

// Bytes defines abstraction for
// formatting gopium struct to byte slice
type Bytes func(gopium.Struct) ([]byte, error)

// Json defines bytes implementation
// which uses json.Marshal with json.Indent to serialize struct
func Json(st gopium.Struct) ([]byte, error) {
	// just use json marshal with indent
	return json.MarshalIndent(st, "", "\t")
}

// Xml defines bytes implementation
// which uses xml.MarshalIndent to serialize struct
func Xml(st gopium.Struct) ([]byte, error) {
	// just use xml marshal with indent
	return xml.MarshalIndent(st, "", "\t")
}

// Csv defines bytes implementation
// that serializes struct to csv format
func Csv(rw io.ReadWriter) Bytes {
	return func(st gopium.Struct) ([]byte, error) {
		// prepare csv writer
		w := csv.NewWriter(rw)
		// write header
		_ = w.Write([]string{
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
		})
		// go through all fields
		// and write then one by one
		for _, f := range st.Fields {
			_ = w.Write([]string{
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
			})
		}
		// flush to buf
		w.Flush()
		// check flush error
		if err := w.Error(); err != nil {
			// this should never happen/covered
			return nil, err
		}
		// and return buf result
		return ioutil.ReadAll(rw)
	}
}
