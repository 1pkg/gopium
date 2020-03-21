package gopium

import "go/types"

// Curator defines system level info curator abstraction
// to expose system word, aligment and cache level sizes
type Curator interface {
	SysWord() int64
	SysAlign() int64
	SysCache(level uint) int64
}

// Exposer defines type info exposer abstraction
// to expose name, size and aligment for provided data type
type Exposer interface {
	Name(types.Type) string
	Size(types.Type) int64
	Align(types.Type) int64
}

// Maven defines Curator abstraction
// and Exposer abstraction aggregation
type Maven interface {
	Curator
	Exposer
}
