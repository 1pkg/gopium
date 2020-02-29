package gopium

import (
	"context"
	"fmt"
	"strings"
)

const (
	VERSION = "1.0.0"
	PKG     = "https://github.com/1pkg/gopium"
	STAMP   = "// 🌺 gopium @1pkg "
	TAG     = "gopium"
)

// Stamped adds gopium stamp
// to specified string
func Stamp(s string) string {
	return fmt.Sprintf("%s - %s", STAMP, s)
}

// Stamped checks if specified string
// has gopium stamp
func Stamped(s string) bool {
	return strings.Contains(s, STAMP)
}

// Runner defines abstraction for
// simple root gopium runner
type Runner interface {
	Run(context.Context) error
}
