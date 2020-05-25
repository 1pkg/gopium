package gopium

import (
	"path"
	"path/filepath"
)

// list of global registered gopium constants
const (
	VERSION = "1.0.0"
	PKG     = "https://github.com/1pkg/gopium"
	STAMP   = "🌺 gopium @1pkg"
)

// Root defines project root path
var Root string

func init() {
	// grabs running root path
	p, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	// until we rich project root
	for path.Base(p) != "gopium" {
		p = path.Dir(p)
	}
	// update the root path
	Root = p
}
