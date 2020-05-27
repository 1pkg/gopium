package gopium

// Categorized defines abstraction for
// categorized structures collection
type Categorized interface {
	Full() map[string]Struct
	Cat(string) (map[string]Struct, bool)
}
