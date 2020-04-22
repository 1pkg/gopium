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
		"not null new reference should return real ref": {
			input: false,
			output: &Reference{
				vals:    make(map[string]int64),
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
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
	// nil ref should alway do default
	r.Set("key", 10)
	val = r.Get("key")
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
	// nil ref should alway do default
	r.Alloc("key")
	val = r.Get("key")
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
	// nil ref should alway do default
	r.Set("key", 10)
	val = r.Get("key")
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
	// nil ref should alway do default
	r.Prune()
	val = r.Get("key")
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
}

func TestRealReferenceMixed(t *testing.T) {
	// stage 0 set up
	var stage int32
	r := NewReference(false)
	defer r.Prune()
	val := r.Get("key")
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
	r.Set("key", 100)
	val = r.Get("key")
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
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
		if val != -1 || atomic.LoadInt32(&stage) != 3 {
			t.Errorf("actual %v doesn't equal to %v", val, -1)
		}
	}()
	// resolved immediately
	go func() {
		val := r.Get("test-100")
		if val != -1 || atomic.LoadInt32(&stage) != 0 {
			t.Errorf("actual %v doesn't equal to %v", val, -1)
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
	if val != -1 {
		t.Errorf("actual %v doesn't equal to %v", val, -1)
	}
}
