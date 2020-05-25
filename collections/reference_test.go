package collections

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewReference(t *testing.T) {
	// prepare
	table := map[string]struct {
		b   bool
		ref *Reference
	}{
		"nil new reference should return nil ref": {
			b:   false,
			ref: nil,
		},
		"not nil new reference should return actual ref": {
			b: true,
			ref: &Reference{
				vals:    make(map[string]interface{}),
				signals: make(map[string]chan struct{}),
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			ref := NewReference(tcase.b)
			// check
			if !reflect.DeepEqual(ref, tcase.ref) {
				t.Errorf("actual %v doesn't equal to %v", ref, tcase.ref)
			}
		})
	}
}

func TestNilReference(t *testing.T) {
	// prepare
	var r *Reference
	r.Set("test-1", 10)
	r.Prune()
	r.Set("test-2", 10)
	r.Alloc("test-3")
	r.Set("test-3", 10)
	table := map[string]struct {
		key string
		val interface{}
	}{
		"invalid key should return empty result": {
			key: "key",
			val: struct{}{},
		},
		"test-1 key should return empty result": {
			key: "test-1",
			val: struct{}{},
		},
		"test-2 key should return empty result": {
			key: "test-2",
			val: struct{}{},
		},
		"test-3 key should return empty result": {
			key: "test-3",
			val: struct{}{},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			val := r.Get(tcase.key)
			// check
			if !reflect.DeepEqual(val, tcase.val) {
				t.Errorf("actual %v doesn't equal to %v", val, tcase.val)
			}
		})
	}
}

func TestActualReference(t *testing.T) {
	// prepare
	var wg sync.WaitGroup
	r := NewReference(true)
	r.Set("test-1", 10)
	go func() {
		r.Set("test-2", 10)
	}()
	r.Alloc("test-3")
	go func() {
		time.Sleep(time.Millisecond)
		r.Set("test-3", 10)
		go func() {
			time.Sleep(time.Millisecond)
			r.Set("test-4", 5)
			time.Sleep(time.Millisecond)
			r.Prune()
		}()
	}()
	r.Alloc("test-4")
	r.Alloc("test-5")
	time.Sleep(3 * time.Millisecond)
	table := map[string]struct {
		key string
		val interface{}
	}{
		"invalid key should return empty result": {
			key: "key",
			val: struct{}{},
		},
		"test-1 key should return empty result": {
			key: "test-1",
			val: struct{}{},
		},
		"test-2 key should return empty result": {
			key: "test-2",
			val: struct{}{},
		},
		"test-3 key should return expected result": {
			key: "test-3",
			val: 10,
		},
		"test-4 key should return expected result": {
			key: "test-4",
			val: 5,
		},
		"test-5 key should return empty result": {
			key: "test-5",
			val: struct{}{},
		},
	}
	for name, tcase := range table {
		// run all parser tests
		// in separate goroutine
		name := name
		tcase := tcase
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Run(name, func(t *testing.T) {
				// exec
				val := r.Get(tcase.key)
				// check
				if !reflect.DeepEqual(val, tcase.val) {
					t.Errorf("actual %v doesn't equal to %v", val, tcase.val)
				}
			})
		}()
	}
	// wait util tests finish
	wg.Wait()
}
