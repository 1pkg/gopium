package mocks

import (
	"io"
	"sync"
)

// Writer defines mock category writer implementation
type Writer struct {
	Gerr  error           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	Cerr  error           `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	RWCs  map[string]*RWC `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	mutex sync.Mutex      `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_     [16]byte        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 64 bytes; struct align: 8 bytes; struct aligned size: 64 bytes; - 🌺 gopium @1pkg

// Generate mock implementation
func (w *Writer) Generate(loc string) (io.WriteCloser, error) {
	// lock rwcs access
	// and init them if they
	// haven't inited before
	defer w.mutex.Unlock()
	w.mutex.Lock()
	if w.RWCs == nil {
		w.RWCs = make(map[string]*RWC)
	}
	// if loc is inside existed rwcs
	// just return found rwc back
	if rwc, ok := w.RWCs[loc]; ok {
		return rwc, w.Gerr
	}
	// otherwise create new rwc
	// store and return it back
	rwc := &RWC{}
	w.RWCs[loc] = rwc
	return rwc, w.Gerr
}

// Category mock implementation
func (w *Writer) Category(cat string) error {
	return w.Cerr
}
