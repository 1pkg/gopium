package gopium

import "go/types"

// Whistleblower defines type info exposer abstraction
// to expose name and size from provided data type
type Whistleblower interface {
	Expose(types.Type) (name string, size int64)
}
