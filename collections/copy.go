package collections

import "github.com/1pkg/gopium/gopium"

// CopyField defines helper that
// deep copies provided field
func CopyField(f gopium.Field) gopium.Field {
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

// CopyStruct defines helper that
// deep copies provided struct
func CopyStruct(s gopium.Struct) gopium.Struct {
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
		ns.Fields = make([]gopium.Field, len(s.Fields), cap(s.Fields))
		copy(ns.Fields, s.Fields)
	}
	// go through struct fields and
	// deep copy them one by one
	for i, f := range ns.Fields {
		ns.Fields[i] = CopyField(f)
	}
	return ns
}
