package mocks

import "bytes"

// RWC defines mock io reader writer closer implementation
type RWC struct {
	buf  bytes.Buffer `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Rerr error        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Werr error        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	Cerr error        `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
	_    [40]byte     `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1,struct_annotate_comment,add_tag_group_force"`
} // struct size: 121 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; - 🌺 gopium @1pkg

// Read mock implementation
func (rwc *RWC) Read(p []byte) (int, error) {
	// in case we have error
	// return it back
	if rwc.Rerr != nil {
		return 0, rwc.Rerr
	}
	// otherwise use buf impl
	return rwc.buf.Read(p)
}

// Write mock implementation
func (rwc *RWC) Write(p []byte) (n int, err error) {
	// in case we have error
	// return it back
	if rwc.Werr != nil {
		return 0, rwc.Werr
	}
	// otherwise use buf impl
	return rwc.buf.Write(p)
}

// Close mock implementation
func (rwc *RWC) Close() error {
	return rwc.Cerr
}
