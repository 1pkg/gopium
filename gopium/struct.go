package gopium

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
