//nolint
//+build pkgs_test_data

package test_walker_visit_top_mixed_level_structs

type FooBar struct {
	xint    int
	xstring string
}

func scope() {
	type FooBarScoped struct {
		xint    int
		xstring string
	}
	func scopescope() {
		type FooBarScopedScoped struct {
			xint    int
			xstring string
		}
	}
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

type FooBarDouble struct {
	f FooBar
	s FooBar
}