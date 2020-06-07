package tests

import (
	"path/filepath"

	"1pkg/gopium"
)

// Gopium defines root gopium path
var Gopium string

// sets gopium data path
func init() {
	// grabs running root path
	p, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	// until we rich project root
	for filepath.Base(p) != gopium.NAME {
		p = filepath.Dir(p)
	}
	Gopium = p
}