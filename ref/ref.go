package ref

import "sync"

// Ref defines backreference helper
// that helps to set, and get wait for
// key value pairs
type Ref struct {
	vals    map[string]int64
	signals map[string]chan struct{}
	mutex   sync.Mutex
}

// NewRef creates ref instance
// accordingly to passed nonnil flag
func NewRef(null bool) *Ref {
	// in case we wanna use
	// nil ref instance
	if null {
		return nil
	}
	// othewise return real ref instance
	return &Ref{
		vals:    make(map[string]int64),
		signals: make(map[string]chan struct{}),
		mutex:   sync.Mutex{},
	}
}

// Get retrieves value for given key
// from the ref in case value hasn't been set yet
// it waits until value will be set
func (r *Ref) Get(key string) int64 {
	// in case of nil ref
	// just skip it and
	// return def size
	if r == nil {
		return -1
	}
	// grab signal with locking
	r.mutex.Lock()
	sig, ok := r.signals[key]
	r.mutex.Unlock()
	// in case there is no slot
	// has been reserved
	// return def size
	if !ok {
		return -1
	}
	// othewise wait for signal
	<-sig
	// lock the ref againg
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// grab the ref value
	if val, ok := r.vals[key]; ok {
		return val
	}
	// in case no value has been set
	// return def size
	return -1
}

// Set update value for given key
// if slot for that value has been preallocated
func (r *Ref) Set(key string, val int64) {
	// in case of nil ref
	// just skip it
	if r == nil {
		return
	}
	// lock the ref
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

// Alloc preallocates slot in the ref
// for the given key
func (r *Ref) Alloc(key string) {
	// in case of nil ref
	// just skip it
	if r == nil {
		return
	}
	// lock the ref
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// if signal hasn't been set yet
	// then allocate a signal for the key
	if _, ok := r.signals[key]; !ok {
		r.signals[key] = make(chan struct{})
	}
}

// Prune releases all value waiters
// and clean all signal resources
func (r *Ref) Prune() {
	// in case of nil ref
	// just skip it
	if r == nil {
		return
	}
	// lock the ref
	defer r.mutex.Unlock()
	r.mutex.Lock()
	// go through all ref's signals
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
