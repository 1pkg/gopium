package gio

import (
	"testing"
)

func TestStdoutMixed(t *testing.T) {
	// prepare
	wc := stdout{}
	// stdout close should be idempotent
	err := wc.Close()
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
	// write should work as expected
	n, err := wc.Write([]byte(""))
	if n != 0 {
		t.Errorf("actual %v doesn't equal to expected %v", n, 0)
	}
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
	// stdout close should be idempotent
	err = wc.Close()
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
	// write should work as expected
	n, err = wc.Write([]byte(""))
	if n != 0 {
		t.Errorf("actual %v doesn't equal to expected %v", n, 0)
	}
	if err != nil {
		t.Errorf("actual %v doesn't equal to expected %v", err, nil)
	}
}
