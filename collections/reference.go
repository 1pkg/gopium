package collections

import "sync"

// Reference defines backreference helper
// that helps to set, and get wait for
// key value pairs
type Reference struct {
	vals    map[string]interface{}
	signals map[string]chan struct{}
	mutex   sync.Mutex
}

// NewReference creates reference instance
// accordingly to passed nil flag
func NewReference(actual bool) *Reference {
	// in case we want to use
	// nil reference instance
	if !actual {
		return nil
	}
	// othewise return actual reference instance
	return &Reference{
		vals:    make(map[string]interface{}),
		signals: make(map[string]chan struct{}),
	}
}

// Get retrieves value for given key from the reference,
// in case value hasn't been set yet
// it waits until value is set
func (r *Reference) Get(key string) interface{} {
	// in case of nil reference
	// just skip it and
	// return def val
	if r == nil {
		return struct{}{}
	}
	// grab signal with locking
	r.mutex.Lock()
	sig, ok := r.signals[key]
	r.mutex.Unlock()
	// in case there is no slot
	// has been reserved
	// return def val
	if !ok {
		return struct{}{}
	}
	// othewise wait for signal
	<-sig
	// lock the reference againg
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// grab the reference value
	if val, ok := r.vals[key]; ok {
		return val
	}
	// in case no value has been set
	// return def val
	return struct{}{}
}

// Set update value for given key,
// if slot for that value has been preallocated
func (r *Reference) Set(key string, val interface{}) {
	// in case of nil reference
	// just skip it
	if r == nil {
		return
	}
	// lock the reference
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// if slot hasn't been allocated yet
	// then just skip set at all
	// otherwise set value for the key
	// and prodcast on the signal
	if ch, ok := r.signals[key]; ok {
		r.vals[key] = val
		// check that channel
		// hasn't been closed yet
		// and then close it
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
}

// Alloc preallocates slot in the
// reference for the given key
func (r *Reference) Alloc(key string) {
	// in case of nil reference
	// just skip it
	if r == nil {
		return
	}
	// lock the reference
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// if signal hasn't been set yet
	// then allocate a signal for the key
	if _, ok := r.signals[key]; !ok {
		r.signals[key] = make(chan struct{})
	}
}

// Prune releases all value waiters
func (r *Reference) Prune() {
	// in case of nil reference
	// just skip it
	if r == nil {
		return
	}
	// lock the reference
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// go through all reference signals
	for _, ch := range r.signals {
		// check that channel
		// hasn't been closed yet
		// and then close it
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
}
