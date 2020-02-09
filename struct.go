package gopium

// Field defines single structure field
// data transfer object abstraction
type Field struct {
	Name     string
	Type     string
	Size     int64
	Tag      string
	Exported bool
	Embedded bool
}

// Struct defines single structure
// data transfer object abstraction
type Struct struct {
	Name   string
	Fields []Field
}

// StructError encapsulates Strategy results
// Struct and error
type StructError struct {
	Struct Struct
	Error  error
}
