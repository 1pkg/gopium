//nolint
//+build pkgs_test_data

package test_walker_visit_top_nested_structs_only

func temp() {
	type fooBar struct {
		xint    int
		xstring string
	}
}

func foo() {
	func bar() {
		type foobar struct {
			xint    int
		}
	}
}