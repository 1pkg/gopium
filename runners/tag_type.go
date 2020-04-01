package runners

// TagType defines type of tag strategy
// that could be run by runners
type TagType string

// list of registered tag types
const (
	None  TagType = "none"
	Soft  TagType = "soft"
	Force TagType = "force"
)
