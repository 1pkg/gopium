//+build tests_data

package multi

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

type ze interface {
	AW() AW
}

type Zeze struct {
	ze
	D
	AZ
	AWA D
}