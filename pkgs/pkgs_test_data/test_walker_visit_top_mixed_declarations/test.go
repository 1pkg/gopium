//nolint
//+build pkgs_test_data

package test_walker_visit_top_mixed_declarations

type FooBarStruct struct {
	structstring string
}

type FooBarCoType FooBarStruct
type FooBarAlias = FooBarStruct
type FooBarStringAlias string

var FooBarAnonymusStruct struct {
	xint int
}

var FooBarFuncVar = func(foobar FooBarStruct) {
	foobar.structstring = "10"
	type FooBarStructNested struct {
		xint int
	}
}

var FooBarAnonymusFunc = func(f struct{ xint int }) error {
	return nil
}

type Test1 struct {
	testuint uint
}

type Test2 struct {
	Test1 uint
}

type Test3 struct {
	Test1 uint
}
