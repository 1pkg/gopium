package tests

import "runtime"

// OnOS acts as ternary operator for os check
// if provided os equals to runtime os
// it returns true value otherwise false value
func OnOS(os string, tval interface{}, fval interface{}) interface{} {
	if runtime.GOOS == os {
		return tval
	}
	return fval
}
