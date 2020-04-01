package gopium

import "fmt"

// Field defines single structure field
// data transfer object abstraction
type Field struct {
	Name     string
	Type     string
	Size     int64
	Align    int64
	Tag      string
	Exported bool
	Embedded bool
	Doc      []string
	Comment  []string
}

// Struct defines single structure
// data transfer object abstraction
type Struct struct {
	Name    string
	Doc     []string
	Comment []string
	Fields  []Field
}

// PadField defines helper that
// creates pad field with specified size
func PadField(pad int64) Field {
	return Field{
		Name:  "_",
		Type:  fmt.Sprintf("[%d]byte", pad),
		Size:  pad,
		Align: 1,
	}
}
