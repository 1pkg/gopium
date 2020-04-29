//+build tests_data

package nested

import "errors"

type A struct {
	a int64
}

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

func scope1() error {
	type B struct {
		b
	}
	type b1 b
	type b2 struct {
		A
		b float64
	}
	return errors.New("test data")
}

func scope2() struct{A complex64, B int64, C float64} {
	// name shadowing 
	type A struct {
		a int32
	}
	type a1 struct {
		i interface{}
	}

	scope3 := func(v int) {
		// name shadowing 
		type a1 struct {
			i struct{}{}
		}
	}

	scope4 := func(v int) {
		// name shadowing 
		var a1 A
		var b1 b
		var c1 C
	}

	return struct{A complex64, B int64, C float64} {}
}

func scope5() interface{} {
	type c1 C
	table := []struct{A string}{{A: "test"}}
	return table[0]
}

type Z struct {
	a bool
	C C
	z bool
}