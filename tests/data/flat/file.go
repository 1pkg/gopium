//+build tests_data

package flat

import (
	"strings"
	"errors"
)

type A struct {
	a int64
}

var a1 string =  strings.Join([]string{"a", "b", "c"}, "|")

type b struct {
	A
	b float64
}

type C struct {
	c []string
	A struct {
		b b
		z A
	}
}

type c1 C
table := []struct{A string}{{A: "test"}}

type D struct {
	t [13]byte
	b bool
	_ int64
}

ggg := func (interface{}){}
type AW func() error

type AZ struct {
	a bool
	D D
	z bool
}