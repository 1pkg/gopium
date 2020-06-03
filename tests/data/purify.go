package data

import (
	"os"
	"strings"

	"1pkg/gopium/tests"
)

// purify helps to transform
// absolute path to relative local one
func purify(loc string) string {
	// remove abs part from loc
	// replace os path separators
	// with underscores and trim them
	loc = strings.Replace(loc, tests.Gopium, "", 1)
	loc = strings.ReplaceAll(loc, string(os.PathSeparator), "_")
	return strings.Trim(loc, "_")
}
