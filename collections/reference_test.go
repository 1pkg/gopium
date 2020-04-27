package collections

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewReference(t *testing.T) {
	// prepare
	table := map[string]struct {
		input  bool
		output *Reference
	}{
		"null new reference should return nil ref": {
			input:  true,
			output: nil,
		},
		"not null new reference should return actual ref": {
			input: false,
			output: &Reference{
				vals:    make(map[string]interface{}),
				signals: make(map[string]chan struct{}),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			output := NewReference(tcase.input)
			// check
			if !reflect.DeepEqual(output, tcase.output) {
				t.Errorf("actual %v doesn't equal to %v", output, tcase.output)
			}
		})
	}
}

func TestNilReferenceMixed(t *testing.T) {
	var r *Reference
	defer r.Prune()
	// nil ref should alway do default
	val := r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
	// nil ref should alway do default
	r.Set("key", 10)
	val = r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
	// nil ref should alway do default
	r.Alloc("key")
	val = r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
	// nil ref should alway do default
	r.Set("key", 10)
	val = r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
	// nil ref should alway do default
	r.Prune()
	val = r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
}

func TestActualReferenceMixed(t *testing.T) {
	// stage 0 set up
	var stage int32
	r := NewReference(false)
	defer r.Prune()
	val := r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
	r.Set("key", 100)
	val = r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
	r.Alloc("key")
	r.Alloc("test-key")
	r.Alloc("test")
	// resolved on set
	go func() {
		val := r.Get("key")
		if val != 10 || atomic.LoadInt32(&stage) != 1 {
			t.Errorf("actual %v doesn't equal to %v", val, 10)
		}
	}()
	// resolved on set
	go func() {
		val := r.Get("key")
		if val != 10 || atomic.LoadInt32(&stage) != 1 {
			t.Errorf("actual %v doesn't equal to %v", val, 10)
		}
	}()
	// resolved on set
	go func() {
		val := r.Get("key")
		if val != 10 || atomic.LoadInt32(&stage) != 1 {
			t.Errorf("actual %v doesn't equal to %v", val, 10)
		}
	}()
	// resolved on update
	go func() {
		val := r.Get("test-key")
		if val != 10 || atomic.LoadInt32(&stage) != 2 {
			t.Errorf("actual %v doesn't equal to %v", val, 10)
		}
		val = r.Get("key")
		if val != 100 || atomic.LoadInt32(&stage) != 2 {
			t.Errorf("actual %v doesn't equal to %v", val, 100)
		}
	}()
	// resolved on prune
	go func() {
		val := r.Get("test")
		if val != struct{}{} || atomic.LoadInt32(&stage) != 3 {
			t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
		}
	}()
	// resolved immediately
	go func() {
		val := r.Get("teststruct{}{}00")
		if val != struct{}{} || atomic.LoadInt32(&stage) != 0 {
			t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
		}
	}()
	// stage 1 set
	time.Sleep(time.Millisecond)
	r.Set("key", 10)
	atomic.AddInt32(&stage, 1)
	val = r.Get("key")
	if val != 10 {
		t.Errorf("actual %v doesn't equal to %v", val, 10)
	}
	// stage 2 update
	time.Sleep(time.Millisecond)
	r.Set("key", 100)
	r.Set("test-key", 10)
	atomic.AddInt32(&stage, 1)
	val = r.Get("key")
	if val != 100 {
		t.Errorf("actual %v doesn't equal to %v", val, 100)
	}
	val = r.Get("test-key")
	if val != 10 {
		t.Errorf("actual %v doesn't equal to %v", val, 10)
	}
	// stage 3 prune
	time.Sleep(time.Millisecond)
	r.Prune()
	atomic.AddInt32(&stage, 1)
	// stage 4 final
	time.Sleep(time.Millisecond)
	atomic.AddInt32(&stage, 1)
	val = r.Get("key")
	if val != struct{}{} {
		t.Errorf("actual %v doesn't equal to %v", val, struct{}{})
	}
}
