//+build tests_data

package multi

import (
	"strings"
)

type A struct {
	a int64
}

var a1 string = strings.Join([]string{"a", "b", "c"}, "|")

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

func scope() {
	type TestAZ struct {
		a bool
		D A
		z bool
	}
}
