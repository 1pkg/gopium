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

// Copy deep copies field
func (f Field) Copy() Field {
	nf := f
	// check that field doc exists
	if f.Doc != nil {
		nf.Doc = make([]string, len(f.Doc), cap(f.Doc))
		copy(nf.Doc, f.Doc)
	}
	// check that field comment exists
	if f.Comment != nil {
		nf.Comment = make([]string, len(f.Comment), cap(f.Comment))
		copy(nf.Comment, f.Comment)
	}
	return nf
}

// Struct defines single structure
// data transfer object abstraction
type Struct struct {
	Name    string
	Doc     []string
	Comment []string
	Fields  []Field
}

// Copy deep copies struct
func (s Struct) Copy() Struct {
	ns := s
	// check that struct doc exists
	if s.Doc != nil {
		ns.Doc = make([]string, len(s.Doc), cap(s.Doc))
		copy(ns.Doc, s.Doc)
	}
	// check that struct comment exists
	if s.Comment != nil {
		ns.Comment = make([]string, len(s.Comment), cap(s.Comment))
		copy(ns.Comment, s.Comment)
	}
	// check that struct fields exists
	if s.Fields != nil {
		ns.Fields = make([]Field, len(s.Fields), cap(s.Fields))
		copy(ns.Fields, s.Fields)
	}
	// go through struct fields and
	// deep copy them one by one
	for i, f := range ns.Fields {
		ns.Fields[i] = f.Copy()
	}
	return ns
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
