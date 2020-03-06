package gopium

import (
	"fmt"
	"strings"
)

const VERSION = "1.0.0"
const PKG = "https://github.com/1pkg/gopium"
const STAMP = " 🌺 gopium @1pkg "

// Stamped adds gopium stamp
// to specified string
func Stamp(s string) string {
	return fmt.Sprintf("%s %s", STAMP, s)
}

// Stamped checks if specified string
// has gopium stamp
func Stamped(s string) bool {
	return strings.Contains(s, STAMP)
}
